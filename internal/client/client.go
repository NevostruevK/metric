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
var serverAddress = "127.0.0.1:8080"

func SetAddress(addr string){
	serverAddress = addr
}

func SendMetrics(sM []metrics.MetricCreater) int {
	client := &http.Client{}
	for i, m := range sM {
		switch obj := m.(type){
		case *metrics.BasicMetric:
			c := clientText{client: client, obj: *obj}
			if err := c.SendMetric(); err!=nil{
				return i
			} 
		case *metrics.Metrics:
			c := clientJSON{client: client, obj: *obj}
			if err := c.SendMetric(); err!=nil{
				return i
			} 
		default:
			fmt.Printf("Type %T not implemented\n", obj)
		}
	}
	return len(sM)
}

type Sender interface{
	SendMetric() error
}

type clientText struct {
	client  *http.Client
	obj 	metrics.BasicMetric
}

type clientJSON struct {
	client  *http.Client
	obj 	metrics.Metrics
}

func (c *clientText) SendMetric() (err error){
	endpoint := url.URL{
		Scheme: "http",
		Host:   serverAddress,
		Path:   "/update/" + c.obj.String(),
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
	if response.StatusCode != http.StatusOK{
		fmt.Printf("SendMetric(Text): read wrong response Status code: %d body %s\n", response.StatusCode,body)
	}
	return nil
}

func (c *clientJSON) SendMetric() (err error){
	endpoint := url.URL{
		Scheme: "http",
		Host:   serverAddress,
		Path:   "/update/",
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
	if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
		body, err = fgzip.Decompress(body)
		if err != nil{
			fmt.Println("SendMetric(JSON): decompress data error: ", err)
			return
		}
	}
	if response.StatusCode != http.StatusOK{
		fmt.Printf("SendMetric(JSON): read wrong response Status code: %d body %s\n", response.StatusCode,body)
	}
	return nil
}
