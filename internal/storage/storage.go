package storage

import (
	"fmt"

	"github.com/NevostruevK/metric/internal/util/metrics"
)
type RepositoryData interface{
	Name() string
	Type() string
	CounterValue() int64
	AddCounterValue(int64) error
}

type Repository interface{
	AddMetric(RepositoryData)
	GetMetric(reqType string, name string) (RepositoryData, error)
	GetAllMetrics() ([]RepositoryData)
}

type MemStorage struct {
	data map[string]RepositoryData
}

func NewMemStorage() *MemStorage{
	return &MemStorage{data: make(map[string]RepositoryData)}
}

func (s *MemStorage) AddMetric(rt RepositoryData) {
	if rt.Type() == metrics.Counter && s.data[rt.Name()] != nil{
		rt.AddCounterValue(s.data[rt.Name()].CounterValue())
	}
	s.data[rt.Name()] = rt
}

func (s *MemStorage) GetMetric(reqType string, name string) (RepositoryData, error){
	if validType := metrics.IsMetricType(reqType); !validType {
		return nil, fmt.Errorf("type %s is not valid metric type",reqType)
	}
	m, ok:= s.data[name]
	if ok{
		if m.Type() == reqType{
			return m, nil
		}
	}
	return nil, fmt.Errorf("type %s : name %s is not valid metric type",reqType,name)
}

func (s *MemStorage) GetAllMetrics() []RepositoryData{
	sM := make([]RepositoryData,0,len(s.data))
	for _, m := range s.data{
		sM = append(sM, m)
	}
	return sM
}

func (s *MemStorage) ShowMetrics(){
	for i, m := range s.data{
		fmt.Println(i, m)
	}
}