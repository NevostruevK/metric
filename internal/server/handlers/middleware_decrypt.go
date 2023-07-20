package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/NevostruevK/metric/internal/util/crypt"
)

func DecryptHanlder(next http.Handler, dcr *crypt.Decrypt) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if dcr == nil {
			next.ServeHTTP(w, r)
			return
		}

		b, err := io.ReadAll(r.Body)
		defer func() {
			err = r.Body.Close()
		}()
		if err != nil {
			msg := fmt.Sprintf("ERROR : failed io.ReadAll with error  %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		data, err := dcr.Decrypt(b)
		if err != nil {
			msg := fmt.Sprintf("ERROR : failed Decrypt with error  %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		next.ServeHTTP(w, r)
	})
}
