package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/NevostruevK/metric/internal/util/fgzip"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

type Agent struct {
	client  *http.Client
	address string
	hashKey string
}

func NewAgent(address, hashKey string) *Agent {
	return &Agent{client: &http.Client{}, address: address, hashKey: hashKey}
}

func SendMetrics(a *Agent, sM []metrics.MetricCreater) int {
	for i, m := range sM {
		switch obj := m.(type) {
		case *metrics.BasicMetric:
			c := clientText{Agent: a, obj: *obj}
			if err := c.SendMetric(); err != nil {
				return i
			}
		case *metrics.Metrics:
			c := clientJSON{Agent: a, obj: *obj}
			if err := c.SendMetric(); err != nil {
				return i
			}
		default:
			fmt.Printf("Type %T not implemented\n", obj)
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
		fmt.Println("SendMetric(Text): create request error: ", err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := c.client.Do(request)
	if err != nil {
		fmt.Println("SendMetric(Text): send request error: ", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("SendMetric(Text): read body error: ", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		fmt.Printf("SendMetric(Text): read wrong response Status code: %d body %s\n", response.StatusCode, body)
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
			return fmt.Errorf(" can't set hash for metric %v , error %v", c.obj, err)
		}
	}
	data, err := json.Marshal(c.obj)
	if err != nil {
		fmt.Println("SendMetric(JSON): marshal data error: ", err)
		return
	}
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("SendMetric(JSON): create request error: ", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		fmt.Println("SendMetric(JSON): send request error: ", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("SendMetric(JSON): read body error: ", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
			body, err = fgzip.Decompress(body)
			if err != nil {
				fmt.Println("SendMetric(JSON): decompress data error: ", err)
				return
			}
		}
		fmt.Printf("SendMetric(JSON): read wrong response Status code: %d body %s\n", response.StatusCode, body)
	}
	return nil
}
