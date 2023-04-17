package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/NevostruevK/metric/internal/util/fgzip"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

type Agent struct {
	client  *http.Client
	address string
	hashKey string
	logger *log.Logger
}

func NewAgent(address, hashKey string) *Agent {
	return &Agent{client: &http.Client{}, address: address, hashKey: hashKey, logger: logger.NewLogger("agent : ",log.LstdFlags|log.Lshortfile)}
}

func SendMetrics(a *Agent, sM []metrics.MetricCreater) int {
	for i, m := range sM {
		switch obj := m.(type) {
		case *metrics.BasicMetric:
			c := clientText{Agent: a, obj: *obj}
			if err := c.SendMetric(); err != nil {
				a.logger.Printf("ERROR : SendMetric returned the error %v\n", err)
				return i
			}
		case *metrics.Metrics:
			c := clientJSON{Agent: a, obj: *obj}
			if err := c.SendMetric(); err != nil {
				a.logger.Printf("ERROR : SendMetric returned the error %v\n", err)
				return i
			}
		default:
			a.logger.Printf("Type %T not implemented\n", obj)
		}
	}
	return len(sM)
}

type Sender interface {
	SendMetric() error
}

type clientText struct {
	*Agent
	obj metrics.BasicMetric
}

type clientJSON struct {
	*Agent
	obj metrics.Metrics
}

func (c *clientText) SendMetric() (err error) {
	endpoint := url.URL{
		Scheme: "http",
		Host:   c.address,
		Path:   "/update/" + c.obj.Type() + "/" + c.obj.Name() + "/" + c.obj.StringValue(),
	}
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), nil)
	if err != nil {
		c.logger.Printf("ERROR : SendMetric(Text): create request error: %v\n", err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := c.client.Do(request)
	if err != nil {
		c.logger.Printf("ERROR : SendMetric(Text): send request error: %v\n", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Printf("ERROR : SendMetric(Text): read body error: %v\n", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		c.logger.Printf("ERROR : SendMetric(Text): read wrong response Status code: %d body %s\n", response.StatusCode, body)
	}
	return nil
}

func (c *clientJSON) SendMetric() (err error) {
	endpoint := url.URL{
		Scheme: "http",
		Host:   c.address,
		Path:   "/update/",
	}

	if c.hashKey != "" {
		if err = c.obj.SetHash(c.hashKey); err != nil {
			msg := fmt.Sprintf("ERROR : SendMetric(JSON): can't set hash for metric %v , error %v\n", c.obj, err)
			c.logger.Println(msg)
			return fmt.Errorf(msg)
		}
	}
	data, err := json.Marshal(c.obj)
	if err != nil {
		c.logger.Printf("ERROR : SendMetric(JSON):json.Marshal error %v\n", err)
		return
	}
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(data))
	if err != nil {
		c.logger.Printf("ERROR : SendMetric(JSON):http.NewRequest error %v\n", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		c.logger.Printf("ERROR : SendMetric(JSON):c.client.Do(request) error %v\n", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Printf("ERROR : SendMetric(JSON):io.ReadAll error %v\n", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
			body, err = fgzip.Decompress(body)
			if err != nil {
				c.logger.Printf("ERROR : SendMetric(JSON):fgzip.Decompress error %v\n", err)
				return
			}
		}
		c.logger.Printf("ERROR : SendMetric(JSON): read wrong response Status code: %d body %s\n", response.StatusCode, body)
	}
	return nil
}
