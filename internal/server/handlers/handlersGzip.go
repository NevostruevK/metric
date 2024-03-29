package handlers

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			msg := fmt.Sprintf("ERROR : CompressHandle:gzip.NewWriterLevel returnen the error : %v", err)
			Logger.Println(msg)
			_, _ = io.WriteString(w, msg)
			return
		}
		defer func() {
			if err = gz.Close(); err != nil {
				Logger.Println(err)
			}
		}()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func DecompressHanlder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			msg := fmt.Sprintf("ERROR : DecompressHanlder:gzip.NewReader returnen the error : %v", err)
			Logger.Println(msg)
			_, _ = io.WriteString(w, msg)
			return
		}
		defer func() {
			if err = gz.Close(); err != nil {
				Logger.Println(err)
			}
		}()
		r.Body = gz
		next.ServeHTTP(w, r)
	})
}
