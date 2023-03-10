package storage

import (
	"fmt"

	"github.com/NevostruevK/metric/internal/util/metrics"
)
type Repository interface{
	AddMetric(m metrics.Metric)
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

func (s *MemStorage) ShowMetrics(){
	for i, m := range s.data{
		fmt.Println(i, m.String())
	}
}