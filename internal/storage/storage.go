package storage

import (
	"fmt"

	"github.com/NevostruevK/metric/internal/util/metrics"
)
type Repository interface{
	AddMetric(m metrics.Metric)
	GetMetric(reqType string, name string) (*metrics.Metric, error)
	GetAllMetrics() ([]metrics.Metric)
}

type MemStorage struct {
	data map[string]metrics.Metric
}

func NewMemStorage() *MemStorage{
	return &MemStorage{data: make(map[string]metrics.Metric)}
}

func (s *MemStorage) AddMetric(m metrics.Metric) {
	if m.Type() == metrics.Counter{
		s.data[m.Name()], _ = m.AddMetricValue(s.data[m.Name()])
		return
	}
	s.data[m.Name()] = m
}
func (s *MemStorage) GetMetric(reqType string, name string) (*metrics.Metric, error){
	if validType := metrics.IsMetricType(reqType); !validType {
		return nil, fmt.Errorf("type %s is not valid metric type",reqType)
	}
	m, ok:= s.data[name]
	if ok{
		if m.Type() == reqType{
			return &m, nil
		}
	}
	return nil, fmt.Errorf("type %s : name %s is not valid metric type",reqType,name)
}

func (s *MemStorage) GetAllMetrics() []metrics.Metric{
	sM := make([]metrics.Metric,0,len(s.data))
	for _, m := range s.data{
		sM = append(sM, m)
	}
	return sM
}

func (s *MemStorage) ShowMetrics(){
	for i, m := range s.data{
		fmt.Println(i, m.String())
	}
}