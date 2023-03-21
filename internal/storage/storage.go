package storage

import (
	"fmt"

	"github.com/NevostruevK/metric/internal/util/metrics"
)
type RepositoryData interface{
	Name() string
	Type() string
//	CounterValue() int64
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
/*	m, ok := rt.(*metrics.Metric)
	if ok && m.Type() == metrics.Counter{
		saved, ok := s.data[rt.Name()].(*metrics.Metric)
		if ok{
			s.data[rt.Name()], _ = m.AddMetricValue(*saved)
			return
		}
	}
*/
	if rt.Type() == metrics.Counter{
		fmt.Printf("New type %T \n",rt)
		fmt.Printf("Old type %T \n",s.data[rt.Name()])
		switch m := rt.(type){
		case *metrics.Metric:
			saved, ok := s.data[rt.Name()].(*metrics.Metric)
			if ok{
				s.data[rt.Name()], _ = m.AddMetricValue(*saved)
				return
			}	
		case metrics.Metrics:
			saved, ok := s.data[rt.Name()].(metrics.Metrics)
			if ok{
				s.data[rt.Name()], _ = m.AddMetricValue(saved)
				return
			}	
		}
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