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
	mJSON := make([]metrics.Metrics, 0, len(sM))
	for i, m := range sM {
		switch obj := m.(type) {
		case *metrics.BasicMetric:
			c := clientText{Agent: a, obj: *obj}
			if err := c.SendMetric(); err != nil {
				return i
			}
		case *metrics.Metrics:
			//			c := clientJSON{Agent: a, obj: *obj}
			mJSON = append(mJSON, *obj)
			//			c := clientJSON{Agent: a, obj: []metrics.Metrics{*obj}}
			//			if err := c.SendMetric(); err != nil {
			//				return i
			//			}
		default:
			fmt.Printf("Type %T not implemented\n", obj)
		}
	}
	if len(mJSON) > 0 {
		c := clientJSON{Agent: a, obj: mJSON}
		if err := c.SendMetric(); err != nil {
			return 0
		}
	}
	return len(sM)
}

/*
type Sender interface {
	SendMetric() error
}
*/

type clientText struct {
	*Agent
	obj metrics.BasicMetric
}

type clientJSON struct {
	*Agent
	obj []metrics.Metrics
}

func (c *clientText) SendMetric() (err error) {
	endpoint := url.URL{
		Scheme: "http",
		Host:   c.address,
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
	if response.StatusCode != http.StatusOK {
		fmt.Printf("SendMetric(Text): read wrong response Status code: %d body %s\n", response.StatusCode, body)
	}
	return nil
}

func (c *clientJSON) SendMetric() (err error) {
	endpoint := url.URL{
		Scheme: "http",
		Host:   c.address,
		//		Path:   "/update/",
	}
	if len(c.obj) > 1 {
		endpoint.Path = "/updates/"
	} else {
		endpoint.Path = "/update/"
	}

	for i, m := range c.obj {
		if c.hashKey != "" {
			if err = c.obj[i].SetHash(c.hashKey); err != nil {
					//				if err = m.SetHash(c.hashKey); err != nil {
					//			if err = c.obj.SetHash(c.hashKey); err != nil {
				return fmt.Errorf(" can't set hash for metric %v , error %v", m, err)
			}
		}
	}
/*	fmt.Println("Agent out inf")
	for _, m := range c.obj {
		fmt.Printf("%s : hash %s\n", m, m.Hash)
	}
*/		//	data, err := json.Marshal(c.obj)
	data, err := json.Marshal(c.obj)
	if err != nil {
		fmt.Println("SendMetric(JSON): marshal data error: ", err)
		return
	}

	data, errCompress := fgzip.Compress(data)

	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("SendMetric(JSON): create request error: ", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	
	if errCompress != nil{
		fmt.Println("can't compress data", err)
	}else{
		fmt.Println("Agent : Content-Encoding : gzip")
		request.Header.Add("Content-Encoding", "gzip")
	}
	
	request.Header.Add("Accept-Encoding", "gzip")

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
//			fmt.Println("Agent Decompress data")
			body, err = fgzip.Decompress(body)
			if err != nil {
				fmt.Println("SendMetric(JSON): decompress data error: ", err)
				return
			}
		}
		fmt.Printf("SendMetric(JSON): read wrong response Status code: %d body %s\n", response.StatusCode, body)
	}

	if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
		fmt.Println("Need to decompress response body")
		body, err = fgzip.Decompress(body)
		if err != nil {
			fmt.Println("SendMetric(JSON): decompress data error: ", err)
			return err
		}
	}
	if !strings.Contains(response.Header.Get("Content-Type"), "application/json") {
		fmt.Println("Agent: response header Content-Type doesn't contain application/json")
		return fmt.Errorf("Agent: response header Content-Type doesn't contain application/json")
	}

	sM := make([]metrics.Metrics, 0, 200)
	
	err = json.Unmarshal(body, &sM)
	if err != nil {
		fmt.Println("Agent: can't unmarshal body")
		return fmt.Errorf("Agent: can't unmarshal body")
	}
	fmt.Println("get data:")
	fmt.Println(sM)
//	fmt.Println("get data:")
//	fmt.Println(body)

	return nil
}
