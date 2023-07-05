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
	"sync/atomic"
	"time"

	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/fgzip"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

const metricsLimit = 1024

type worker struct {
	client  *http.Client
	address string
	hashKey string
	logger  *log.Logger
}

type workers struct {
	free    int32
	workers []worker
}

func NewWorker(address, hashKey string, id int) *worker {
	name := fmt.Sprintf("worker %d ", id)
	return &worker{client: &http.Client{}, address: address, hashKey: hashKey, logger: logger.NewLogger(name, log.LstdFlags|log.Lshortfile)}
}

func (w *worker) start(ctx context.Context, inCh, reuseCh chan []metrics.Metrics, free *int32) {
	w.logger.Println("Start")
	atomic.AddInt32(free, 1)
	for {
		select {
		case newM := <-inCh:
			w.logger.Printf("Get %d metrics", len(newM))
			atomic.AddInt32(free, -1)
			send, err := w.Send(ctx, newM)
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
			w.logger.Println("Finished")
			atomic.AddInt32(free, 1)
			return
		}
	}
}

func (w *worker) Send(ctx context.Context, sM []metrics.Metrics) (int, error) {
	endpoint := url.URL{
		Scheme: "http",
		Host:   w.address,
		Path:   "/update/",
	}
	for i, m := range sM {
		if w.hashKey != "" {
			if err := m.SetHash(w.hashKey); err != nil {
				msg := fmt.Sprintf("ERROR : SendMetric(JSON): can't set hash for metric %v , error %v\n", m, err)
				w.logger.Println(msg)
				return i, fmt.Errorf(msg)
			}
		}
		data, err := json.Marshal(m)
		if err != nil {
			w.logger.Printf("ERROR : SendMetric(JSON):json.Marshal error %v\n", err)
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
		defer func() {
			err = response.Body.Close()
		}()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			w.logger.Printf("ERROR : SendMetric(JSON):io.ReadAll error %v\n", err)
			return i, err
		}
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

func StartAgent(ctx context.Context, cmd *commands.Commands) {
	lgr := logger.NewLogger("agent : ", log.LstdFlags|log.Lshortfile)
	lgr.Println("Start")
	chIn := make(chan []metrics.Metrics)
	chOut := make(chan []metrics.Metrics)
	w := &workers{free: 0, workers: make([]worker, 0, cmd.RateLimit)}
	for i := 0; i < cmd.RateLimit; i++ {
		w.workers = append(w.workers, *NewWorker(cmd.Address, cmd.Key, i))
		go w.workers[i].start(ctx, chOut, chIn, &w.free)
	}
	go CollectMetrics(ctx, cmd.PollInterval, chIn)
	reportTicker := time.NewTicker(cmd.ReportInterval)
	defer reportTicker.Stop()
	sM := make([]metrics.Metrics, 0, metricsLimit)

	for {
		select {
		case newM := <-chIn:
			sM = append(sM, newM...)
			lgr.Printf("Recieve: %d metrics", len(sM))
		case <-reportTicker.C:
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
		case <-ctx.Done():
			lgr.Println("Finished")
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
