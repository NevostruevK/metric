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
//var client = nil


func SendMetrics(sM []metrics.Metric, size int){
    for i, m := range sM{
//        fmt.Println(i," : ",m.String())
        SendMetric(m)
        if (i>=(size-1)){
             break
        }     
//        fmt.Println("---------------- ")
    }
}
    
func SendMetric(sM metrics.Metric) {
//func SendMetric(sM metrics.Metric) {
		//sM := metrics.Get()
    // адрес сервиса (как его писать, расскажем в следующем уроке)
    endpoint := serverAddress + "update/"
//    endpoint := "http://127.0.0.1:8080/"
    // контейнер данных для запроса
//    data := url.Values{}
    // приглашение в консоли
    // конструируем HTTP-клиент
    client := &http.Client{}
//	http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>;
//	Content-Type: text/plain.
	
//	data.Set("type", sM.GetType())
//	data.Add("name", sM.GetName())
//	data.Add("value", sM.ValueToString())
// конструируем запрос
    // запрос методом POST должен, кроме заголовков, содержать тело
    // тело должно быть источником потокового чтения io.Reader
    // в большинстве случаев отлично подходит bytes.Buffer
/*    request, err := http.NewRequest(http.MethodPost, 
        endpoint + sM.String(), 
        bytes.NewBufferString(data.Encode()))
*/
    
    request, err := http.NewRequest(http.MethodPost, endpoint + sM.String(), nil)
    if err != nil {
        fmt.Println("http.NewRequest",err)
        os.Exit(1)
    }
    // в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
    request.Header.Add("Content-Type", "text/plain")
//    request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
    // отправляем запрос и получаем ответ
//    fmt.Println("Server send request :",endpoint + sM.String())
            // создаём новый Recorder
/*    w := httptest.NewRecorder()
            // определяем хендлер
    h := http.HandlerFunc(handlers.URLHandler)
            // запускаем сервер
    h.ServeHTTP(w, request)
*/
//    response := w.Result()

    response, err := client.Do(request)
//    _, err = client.Do(request)
    if err != nil {
        fmt.Println("client.Do",err)
        os.Exit(1)
    }

//    fmt.Println("Server response :",response.Status)
    // печатаем код ответа
//    fmt.Println("Статус-код ", response.Status)
    defer response.Body.Close()
    // читаем поток из тела ответа
//    body, err := io.ReadAll(response.Body)
    body, err := io.ReadAll(response.Body)
    if err != nil {
        fmt.Println("io.ReadAll",err)
        os.Exit(1)
    }
    // и печатаем его
    fmt.Println("BODY: ",string(body))




}