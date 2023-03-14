package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/NevostruevK/metric/internal/server"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

type client struct {
	*http.Client
}

func SendMetrics(sM []metrics.Metric) {
		c := &client{&http.Client{}}

	for _, m := range sM {
		c.SendMetric(m)
	}
}

func (c *client) SendMetric(sM metrics.Metric) {
	endpoint := url.URL{
		Scheme: "http",
		Host:   server.ServerAddress,
		Path:   "/update/" + sM.String(),
	}
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), nil)
	if err != nil {
		fmt.Println("http.NewRequest", err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := c.Do(request)
	if err != nil {
		fmt.Println("Send request error", err)
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
}
