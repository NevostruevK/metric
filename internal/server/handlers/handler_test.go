package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NevostruevK/metric/internal/server/handlers"
	"github.com/NevostruevK/metric/internal/storage"
)

func TestSaveMetricHandler(t *testing.T) {
	s := storage.NewMemStorage()

type want struct {
        code        int
    }
	tests := []struct {
		name string
		path string
		want want
	}{
		{
            name: "positive test with counter #1",
			path: "/update/counter/nameCounter/100",
            want: want{
                code:        200,
            },
        },
		{
            name: "negative test with a wrong value",
			path: "/update/gauge/testGauge/qwert",
            want: want{
                code:        400,
            },
        },
		{
            name: "negative test with a wrong type",
			path: "/update/myType/testGauge/0.123",
            want: want{
                code:        501,
            },
        },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
				request := httptest.NewRequest(http.MethodPost, tt.path, nil)

				w := httptest.NewRecorder()
				h := http.HandlerFunc(handlers.AddMetricHandler(s))
				h.ServeHTTP(w, request)
				res := w.Result()
	
				if res.StatusCode != tt.want.code {
					t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
				}
				defer res.Body.Close()
		})
	}
}

