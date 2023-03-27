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
		case *metrics.Metric:
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
	obj 	metrics.Metric
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
		fmt.Println("http.NewRequest", err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := c.client.Do(request)
	if err != nil {
		fmt.Println(c.obj," : Send request error", err)
		return
	}
	fmt.Println("response Status code : ", response.StatusCode)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("io.ReadAll", err)
		return
	}
	fmt.Println("response body: ", string(body))
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
		fmt.Println("json.Marshal", err)
		return
	}
//	data, errCompress := fgzip.Compress(data)

	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("http.NewRequest", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

/*	if errCompress != nil{
		fmt.Println("can't compress data", err)
	}else{
		request.Header.Add("Content-Encoding", "gzip")
	}
*/
	response, err := c.client.Do(request)
	if err != nil {
		fmt.Println("Send request error", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("io.ReadAll", err)
		return
	}
	if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
//		if response.Header.Get("Content-Encoding") == ("gzip"){
		body, err = fgzip.Decompress(body)
		if err != nil{
			fmt.Println("can't decompress data", err)
			return
		}
	}
//	body, err := decompressGzip(response.Body)

	if response.StatusCode != http.StatusOK{
		fmt.Println("response Status code : ", response.StatusCode)
		fmt.Println("response body: ", string(body))
	}
	return nil
}
