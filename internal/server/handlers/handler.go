package handlers

import (
	"fmt"
	"net/http"
)

func URLHandler(w http.ResponseWriter, r *http.Request) {
        // извлекаем фрагмент query= из URL запроса search?query=something
        //    q := r.URL
        //    fmt.Println("request URL",q)
        //      ct:= r.Header.Get("Content-Type")
        //    fmt.Println("request Header",ct)
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "Server response")
}