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
	saver *saver
	needToSyncWrite bool
}

func NewMemStorage(restore bool, needToSyncWrite bool, filename string) *MemStorage{
	data := make(map[string]RepositoryData)
	if filename == ""{
		return &MemStorage{data: data, saver: nil, needToSyncWrite: false}		
	}
	s, err := NewSaver(filename)
	if err!=nil{
		fmt.Printf("Can't write metrics to %s\n",filename)
		return &MemStorage{data: data, saver: nil, needToSyncWrite: false}	
	}
	if restore{
		l, err := NewLoader(filename)
		if err!=nil{
			fmt.Printf("Can't load metrics from %s\n",filename)
		}else{
			defer l.Close()
			for {
				m, err := l.ReadMetric()
				if err!=nil{
					break
				}
				data[m.Name()] = m
			}
		}
	}
	return &MemStorage{data: data, saver: s, needToSyncWrite: needToSyncWrite,}
}
func (s *MemStorage) SaveAllIntoFile() (int,error){
	if s.saver == nil{
		return 0, fmt.Errorf("can't save metrics into file, saver wasn't initialized")
	}
	count := 0
	for _, m := range s.data{
		if err := s.saver.WriteMetric(m); err!=nil{
			return count, fmt.Errorf("can't save metric into file, encoder error")
		}
		count++
	}
	return count, nil
}

func (s *MemStorage) Close(){
	s.saver.Close()
}

func (s *MemStorage) AddMetric(rt RepositoryData) {
	if rt.Type() == metrics.Counter && s.data[rt.Name()] != nil{
		rt.AddCounterValue(s.data[rt.Name()].CounterValue())
	}
	s.data[rt.Name()] = rt
	if (s.needToSyncWrite){
		s.saver.WriteMetric(rt)
	}
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