package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

// GetMetricJSONHandler обработчик запроса /value/.
// Возвращает метрику, полученную в теле запроса.
func GetMetricJSONHandler(s storage.Repository, hashKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sM, code, err := getMetricFromRequest(r, "", 1)
		if err != nil {
			msg := fmt.Sprintf(" ERROR : GetMetricJSONHandler:getMetricFromRequest returned the error : %v", err)
			Logger.Println(msg)
			http.Error(w, msg, code)
			return
		}
		m := sM[0]
		rt, err := s.GetMetric(context.Background(), m.MType, m.ID)
		if err != nil {
			msg := fmt.Sprintf(" ERROR : GetMetricJSONHandler:GetMetric returned the error : %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusNotFound)
			return
		}
		m = rt.ConvertToMetrics()

		if code, err = sendResponse(w, []metrics.Metrics{m}, false, hashKey); err != nil {
			msg := fmt.Sprintf(" ERROR : GetMetricJSONHandler:sendResponse returned the error : %v", err)
			Logger.Println(msg)
			http.Error(w, msg, code)
		}
	}
}

func prepareMetricsForStorage(sM []metrics.Metrics) ([]metrics.Metrics, error) {
	st := storage.NewMemStorage(false, false, "")
	if err := st.AddGroupOfMetrics(context.Background(), sM); err != nil {
		return nil, err
	}
	pM, err := st.GetAllMetrics(context.Background())
	if err != nil {
		return nil, err
	}
	return pM, nil
}

// AddBatchMetricJSONHandler обработчик запроса /updates/.
// Принимает массив метри из тела запроса.
func AddBatchMetricJSONHandler(s storage.Repository, hashKey string, cap int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sM, code, err := getMetricFromRequest(r, hashKey, cap)
		if err != nil {
			msg := fmt.Sprintf(" ERROR : AddBatchMetricJSONHandler:getMetricFromRequest returned the error : %v", err)
			Logger.Println(msg)
			http.Error(w, msg, code)
			return
		}

		sM, err = prepareMetricsForStorage(sM)
		if err != nil {
			msg := fmt.Sprintf(" ERROR : AddBatchMetricJSONHandler:prepareMetricsForStorage returned the error : %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if err = s.AddGroupOfMetrics(context.Background(), sM); err != nil {
			msg := fmt.Sprintf(" ERROR : AddBatchMetricJSONHandler:AddGroupOfMetrics returned the error : %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if code, err = sendResponse(w, sM, true, hashKey); err != nil {
			msg := fmt.Sprintf(" ERROR : AddBatchMetricJSONHandler:sendResponse returned the error : %v", err)
			Logger.Println(msg)
			http.Error(w, msg, code)
		}
	}
}

// AddMetricJSONHandler обработчик запроса /update/.
// Принимает метрику из тела запроса.
func AddMetricJSONHandler(s storage.Repository, hashKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sM, code, err := getMetricFromRequest(r, hashKey, 1)
		if err != nil {
			msg := fmt.Sprintf(" ERROR : AddMetricJSONHandler:getMetricFromRequest returned the error : %v", err)
			Logger.Println(msg)
			http.Error(w, msg, code)
			return
		}
		if err = s.AddMetric(context.Background(), &sM[0]); err != nil {
			msg := fmt.Sprintf(" ERROR : AddMetricJSONHandler:AddMetric returned the error : %v", err)
			Logger.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if code, err = sendResponse(w, sM, false, hashKey); err != nil {
			msg := fmt.Sprintf(" ERROR : AddMetricJSONHandler:sendResponse returned the error : %v", err)
			Logger.Println(msg)
			http.Error(w, msg, code)
		}
	}
}

func getMetricFromRequest(r *http.Request, hashKey string, initialCapacity int) ([]metrics.Metrics, int, error) {

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, http.StatusBadRequest, fmt.Errorf("Content-Type is not application/json")
	}

	b, err := io.ReadAll(r.Body)
	defer func() {
		err = r.Body.Close()
	}()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	sM := make([]metrics.Metrics, 0, initialCapacity)

	if r.URL.Path != "/updates/" {
		m := metrics.Metrics{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		sM = append(sM, m)
	} else {
		err = json.Unmarshal(b, &sM)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	for _, m := range sM {
		if m.MType == "" || m.ID == "" {
			return nil, http.StatusBadRequest, fmt.Errorf("param is missed : id %s, type %s", m.ID, m.MType)
		}

		if isValidType := metrics.IsMetricType(m.MType); !isValidType {
			return nil, http.StatusNotImplemented, fmt.Errorf("type %s is not implemented", m.MType)
		}

		if hashKey != "" {
			ok, err := m.CheckHash(hashKey)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
			if !ok {
				return nil, http.StatusBadRequest, fmt.Errorf("wrong hash for metric %s", m.String())
			}
		}
	}

	return sM, http.StatusOK, nil
}

func sendResponse(w http.ResponseWriter, sM []metrics.Metrics, sendSlice bool, hashKey string) (int, error) {

	for i, m := range sM {
		if hashKey != "" {
			if err := sM[i].SetHash(hashKey); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("%w for metric %s", err, m.String())
			}
		}
	}
	var data []byte
	var err error
	if sendSlice {
		data, err = json.Marshal(&sM)
	} else {
		data, err = json.Marshal(sM[0])
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
	return http.StatusOK, nil
}
