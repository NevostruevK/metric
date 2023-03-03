package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/NevostruevK/metric/internal/util/metrics"
)

const SERVER_ADDRES = "http://localhost:8080/"
//const SERVER_ADDRES = "http://127.0.0.1:8080/"

//func SendMetric(sM []metrics.Metric) {
func SendMetric(sM metrics.Metric) {
		//sM := metrics.Get()
    // адрес сервиса (как его писать, расскажем в следующем уроке)
    endpoint := SERVER_ADDRES
//    endpoint := "http://127.0.0.1:8080/"
    // контейнер данных для запроса
    data := url.Values{}
    // приглашение в консоли
    
    // конструируем HTTP-клиент
    client := &http.Client{}
//	http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>;
//	Content-Type: text/plain.
	
	data.Set("type", sM.GetType())
	data.Add("name", sM.GetName())
	data.Add("value", sM.GetValue())
// конструируем запрос
    // запрос методом POST должен, кроме заголовков, содержать тело
    // тело должно быть источником потокового чтения io.Reader
    // в большинстве случаев отлично подходит bytes.Buffer
    request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(data.Encode()))
    if err != nil {
        fmt.Println("http.NewRequest",err)
        os.Exit(1)
    }
    // в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
    request.Header.Add("Content-Type", "text/plain")
    request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
    // отправляем запрос и получаем ответ
    response, err := client.Do(request)
    if err != nil {
        fmt.Println("client.Do",err)
        os.Exit(1)
    }
    // печатаем код ответа
    fmt.Println("Статус-код ", response.Status)
    defer response.Body.Close()
    // читаем поток из тела ответа
    body, err := io.ReadAll(response.Body)
    if err != nil {
        fmt.Println("io.ReadAll",err)
        os.Exit(1)
    }
    // и печатаем его
    fmt.Println("BODY: ",string(body))
}