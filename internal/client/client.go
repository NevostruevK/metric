// Package client Клиент для сервера по сбору метрик.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/NevostruevK/metric/proto"

	grpcclient "github.com/NevostruevK/metric/internal/grpc/client"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/crypt"
	"github.com/NevostruevK/metric/internal/util/fgzip"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

const metricsLimit = 1024

//const timeOutForSending = time.Second

type worker struct {
	client *http.Client
	gRPC   *grpcclient.Client
	cfg    *commands.Config
	crypt  *crypt.Crypt
	logger *log.Logger
}

type workers struct {
	free    int32
	workers []worker
}

func NewWorker(cfg *commands.Config, id int, crypt *crypt.Crypt) (*worker, error) {
	name := fmt.Sprintf("worker %d ", id)
	if !cfg.GRPC {
		return &worker{
			client: &http.Client{},
			cfg:    cfg,
			logger: logger.NewLogger(name, log.LstdFlags|log.Lshortfile),
			crypt:  crypt,
		}, nil
	}
	client, err := grpcclient.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &worker{
		gRPC:   client,
		cfg:    cfg,
		logger: logger.NewLogger(name, log.LstdFlags|log.Lshortfile),
	}, nil

}
func (w *worker) Close() error {
	if w.gRPC != nil {
		return w.gRPC.Close()
	}
	return nil
}

func (w *worker) start(ctx context.Context, inCh, reuseCh chan []metrics.Metrics, free *int32, wg *sync.WaitGroup) {
	defer wg.Done()
	w.logger.Println("Start")
	atomic.AddInt32(free, 1)
	var (
		send int
		err  error
	)
	for {
		select {
		case newM := <-inCh:
			w.logger.Printf("Get %d metrics", len(newM))
			atomic.AddInt32(free, -1)
			if w.gRPC != nil {
				send, err = w.SendGRPC(context.Background(), newM)
			} else {
				send, err = w.Send(context.Background(), newM)
			}
			if err != nil {
				w.logger.Printf("Error : %v", err)
			}
			w.logger.Printf("Sended %d from %d metrics", send, len(newM))
			if send < len(newM) {
				newM = newM[send:]
				reuseCh <- newM
			}
			atomic.AddInt32(free, 1)
			w.logger.Println("Work compleate")
		case <-ctx.Done():
			if err = w.Close(); err != nil {
				w.logger.Printf("Finished with %v", err)
			} else {
				w.logger.Println("Finished normal")
			}
			atomic.AddInt32(free, 1)
			return
		}
	}
}

func PrepareDataForMetric(m metrics.Metrics, hashKey string, crypt *crypt.Crypt) ([]byte, error) {
	if hashKey != "" {
		if err := m.SetHash(hashKey); err != nil {
			return nil, err
		}
	}
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	if crypt != nil {
		return crypt.Crypt(data)
	}
	return data, nil
}

func (w *worker) SendGRPC(ctx context.Context, sM []metrics.Metrics) (int, error) {
	op := "SendGRPC: "
	if w.cfg.HashKey != "" {
		for i, m := range sM {
			if err := sM[i].SetHash(w.cfg.HashKey); err != nil {
				return 0, fmt.Errorf("%s for metric %v failed: %v", op, m, err)
			}
		}
	}
	for i, m := range sM {
		_, err := w.gRPC.C.AddMetric(
			ctx,
			&pb.AddMetricRequest{Metric: m.ToProto()},
		)
		if err != nil {
			return i, fmt.Errorf("%s for metric %v failed: %v", op, m, err)
		}
	}
	return len(sM), nil
}

func (w *worker) Send(ctx context.Context, sM []metrics.Metrics) (int, error) {
	endpoint := url.URL{
		Scheme: "http",
		Host:   w.cfg.Address,
		Path:   "/update/",
	}
	for i, m := range sM {
		data, err := PrepareDataForMetric(m, w.cfg.HashKey, w.crypt)
		if err != nil {
			w.logger.Printf("ERROR : PrepareDataForMetric %v failed with %v\n", m, err)
			return i, err
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewBuffer(data))
		if err != nil {
			w.logger.Printf("ERROR : SendMetric(JSON):http.NewRequest error %v\n", err)
			return i, err
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := w.client.Do(request)
		if err != nil {
			w.logger.Printf("ERROR : SendMetric(JSON):c.client.Do(request) error %v\n", err)
			return i, err
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			w.logger.Printf("ERROR : SendMetric(JSON):io.ReadAll error %v\n", err)
			return i, err
		}
		response.Body.Close()
		if response.StatusCode != http.StatusOK {
			if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
				body, err = fgzip.Decompress(body)
				if err != nil {
					w.logger.Printf("ERROR : SendMetric(JSON):fgzip.Decompress error %v\n", err)
					return i, err
				}
			}
			w.logger.Printf("ERROR : SendMetric(JSON): read wrong response Status code: %d body %s\n", response.StatusCode, body)
		}
	}
	return len(sM), nil
}

func StartAgent(ctx context.Context, cmd *commands.Config, complete chan struct{}) {
	lgr := logger.NewLogger("agent : ", log.LstdFlags|log.Lshortfile)
	lgr.Println("Start")
	chIn := make(chan []metrics.Metrics)
	chOut := make(chan []metrics.Metrics)
	w := &workers{free: 0, workers: make([]worker, 0, cmd.RateLimit)}
	cr, err := crypt.NewCrypt(cmd.CryptoKey)
	if err != nil {
		lgr.Printf("failed to create crypt entity %v", err)
	}

	wctx, wcancel := context.WithCancel(context.Background())
	defer wcancel()
	wg := &sync.WaitGroup{}

	for i := 0; i < cmd.RateLimit; i++ {
		newWorker, err := NewWorker(cmd, i, cr)
		if err != nil {
			lgr.Fatalf("StartAgent failed: %v", err)
		}
		w.workers = append(w.workers, *newWorker)
		wg.Add(1)
		go w.workers[i].start(wctx, chOut, chIn, &w.free, wg)
	}
	go CollectMetrics(ctx, time.Duration(cmd.PollInterval.Duration), chIn)
	reportTicker := time.NewTicker(cmd.ReportInterval.Duration)
	defer reportTicker.Stop()
	sM := make([]metrics.Metrics, 0, metricsLimit)

	giveJob := func() {
		free := int(atomic.LoadInt32(&w.free))
		lgr.Printf("I have %d workers for %d metrics", free, len(sM))
		for i := free; i > 0; i-- {
			size := len(sM) / i
			if size == 0 {
				continue
			}
			chOut <- sM[:size]
			sM = sM[size:]
		}
		lgr.Println("Gave Job")
		sM = nil
	}

	for {
		select {
		case newM := <-chIn:
			sM = append(sM, newM...)
			lgr.Printf("Recieve: %d metrics", len(sM))
		case <-reportTicker.C:
			giveJob()
		case <-ctx.Done():
			lgr.Println("receive signal for finishing")
			giveJob()
			wcancel()
			lgr.Println("wait for finishing workers")
			wg.Wait()
			lgr.Println("workers finished")
			complete <- struct{}{}
			return
		}
	}

}

func CollectMetrics(ctx context.Context, pollInterval time.Duration, ch chan []metrics.Metrics) {
	lgr := logger.NewLogger("collect : ", log.LstdFlags|log.Lshortfile)

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			sM, _ := metrics.GetAdvanced()
			sM = append(sM, metrics.Get()...)
			lgr.Printf("Send: %d metrics", len(sM))
			ch <- sM
		case <-ctx.Done():
			return
		}
	}
}
