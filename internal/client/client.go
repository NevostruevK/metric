package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/NevostruevK/metric/internal/server"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

type client struct{
        *http.Client
}

func SendMetrics(sM []metrics.Metric, size int) {
        c := &client{&http.Client{}}

        for i, m := range sM {
                c.SendMetric(m)
                if i >= (size - 1) {
                        break
                }
        }
}

func(c *client) SendMetric(sM metrics.Metric) {
        endpoint := url.URL{
                Scheme: "http",
                Host: server.ServerAddress,
                Path: "/update/" + sM.String(),
        }
        request, err := http.NewRequest(http.MethodPost, endpoint.String(), nil)
                if err != nil {
                fmt.Println("http.NewRequest", err)
                os.Exit(1)
        }
        request.Header.Set("Content-Type", "text/plain")
        //    request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
        //    fmt.Println("Server send request :",endpoint + sM.String())
        /*           //  создаём новый Recorder
        w := httptest.NewRecorder()
                // определяем хендлер
        h := http.HandlerFunc(handlers.URLHandler)
                // запускаем сервер
        h.ServeHTTP(w, request)
        response := w.Result()
        */
        response, err := c.Do(request)
        if err != nil {
                fmt.Println("Send request error", err)
                os.Exit(1)
        }
        fmt.Println("response Status code : " ,response.StatusCode)
        defer response.Body.Close()
        body, err := io.ReadAll(response.Body)
        if err != nil {
                fmt.Println("io.ReadAll", err)
                os.Exit(1)
        }
        fmt.Println("response body: ", string(body))
}