package handlers

import (
	"log"
	"net/http"
)

func LoggerHanlder(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println("URL : ", r.URL)
		next.ServeHTTP(w, r)
	})
}
