package client

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/NevostruevK/metric/internal/util/metrics"
)

//const SERVER_ADDRES = "http://localhost:8080/"
const serverAddress = "http://127.0.0.1:8080/"

func SendMetrics(sM []metrics.Metric, size int) {
        for i, m := range sM {
                SendMetric(m)
                if i >= (size - 1) {
                        break
                }
        }
}

func SendMetric(sM metrics.Metric) {
        endpoint := serverAddress + "update/"
        client := &http.Client{}

        request, err := http.NewRequest(http.MethodPost, endpoint+sM.String(), nil)
        if err != nil {
                fmt.Println("http.NewRequest", err)
                os.Exit(1)
        }
        request.Header.Add("Content-Type", "text/plain")
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

        response, err := client.Do(request)
        if err != nil {
                fmt.Println("client.Do", err)
                os.Exit(1)
        }

        defer response.Body.Close()
        body, err := io.ReadAll(response.Body)
        if err != nil {
                fmt.Println("io.ReadAll", err)
                os.Exit(1)
        }
        fmt.Println("BODY: ", string(body))
}