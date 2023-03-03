package main

import (
	"fmt"

	"github.com/NevostruevK/metric/internal/client"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

const SERVER_ADDRES = "http://localhost:8080/"
//const SERVER_ADDRES = "http://127.0.0.1:8080/"

func main() {
	fmt.Println("sM := metrics.Get()")
	sM := metrics.Get()
	fmt.Println("client.SendMetric(sM)")
//	client.SendMetric(sM)
	client.SendMetric(sM[0])
}