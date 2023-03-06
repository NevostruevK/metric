package main

import "fmt"

func main() {
	fmt.Printf("What the Fuck? %.2f \n", 0.123)
}

/*
import (
	"log"
	"net/http"

	// важно: путь до пакета handlers в вашем случае может быть другим
	"github.com/NevostruevK/metric/internal/server/handlers"
	//	"../../internal/server/handlers"
	//	"d:Golang/metrics/metric/internal/server/handlers"
)

func main() {
	http.HandleFunc("/status", handlers.StatusHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
} */