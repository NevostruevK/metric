package storage

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

type RepositoryData interface {
	Name() string
	Type() string
	StringValue() string
	CounterValue() int64
	GaugeValue() float64
	AddCounterValue(int64) error
	ConvertToMetrics() metrics.Metrics
}

type Repository interface {
	AddMetric(context.Context, RepositoryData) error
	GetMetric(ctx context.Context, reqType, name string) (RepositoryData, error)
	GetAllMetrics(context.Context) ([]metrics.Metrics, error)
	AddGroupOfMetrics(ctx context.Context, sM []metrics.Metrics) error
	Ping() error
}

type MemStorage struct {
	data            map[string]RepositoryData
	saver           *saver
	logger          *log.Logger
	needToSyncWrite bool
	mu              sync.RWMutex
}

func NewMemStorage(restore, needToSyncWrite bool, filename string) *MemStorage {
	lgr := logger.NewLogger("mem storage : ", log.LstdFlags|log.Lshortfile)
	data := make(map[string]RepositoryData)
	s, err := NewSaver(filename)
	if err != nil {
		lgr.Printf("Can't write metrics to %s\n", filename)
		return &MemStorage{data: data, saver: s, logger: lgr, needToSyncWrite: false}
	}
	if filename == "" {
		return &MemStorage{data: data, saver: s, logger: lgr, needToSyncWrite: false}
	}
	if restore {
		l, err := NewLoader(filename)
		if err != nil {
			lgr.Printf("Can't load metrics from %s\n", filename)
		} else {
			defer l.Close()
			for {
				m, err := l.ReadMetric()
				if err != nil {
					break
				}
				data[m.Name()] = m
			}
		}
	}
	return &MemStorage{data: data, saver: s, logger: lgr, needToSyncWrite: needToSyncWrite}
}
func (s *MemStorage) SaveAllIntoFile() (int, error) {
	if s.saver == nil {
		s.logger.Println("can't save metrics into file, saver wasn't initialized")
		return 0, fmt.Errorf("can't save metrics into file, saver wasn't initialized")
	}
	count := 0
	for _, m := range s.data {
		if err := s.saver.WriteMetric(m); err != nil {
			msg := fmt.Sprintf("ERROR : can't save metric into file, encoder error %v\n", err)
			s.logger.Println(msg)
			return count, fmt.Errorf(msg)
		}
		count++
	}
	return count, nil
}

func (s *MemStorage) Close() error {
	return s.saver.Close()
}

func (s *MemStorage) Ping() error {
	return fmt.Errorf("can't ping memory storage")
}

func (s *MemStorage) AddGroupOfMetrics(ctx context.Context, sM []metrics.Metrics) error {
	for i, m := range sM {
		if err := s.AddMetric(ctx, &sM[i]); err != nil {
			return fmt.Errorf("can't AddMetric %s", m)
		}
	}
	return nil
}

func (s *MemStorage) AddMetric(ctx context.Context, rt RepositoryData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if rt.Type() == metrics.Counter && s.data[rt.Name()] != nil {
		rt.AddCounterValue(s.data[rt.Name()].CounterValue())
	}
	s.data[rt.Name()] = rt
	if s.needToSyncWrite {
		s.saver.WriteMetric(rt)
	}
	return nil
}

func (s *MemStorage) GetMetric(ctx context.Context, reqType, name string) (RepositoryData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if validType := metrics.IsMetricType(reqType); !validType {
		return nil, fmt.Errorf("type %s is not valid metric type", reqType)
	}
	m, ok := s.data[name]
	if ok {
		if m.Type() == reqType {
			return m, nil
		}
	}
	return nil, fmt.Errorf("type %s : name %s is not valid metric type", reqType, name)
}

func (s *MemStorage) GetAllMetrics(context.Context) ([]metrics.Metrics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sM := make([]metrics.Metrics, 0, len(s.data))
	for _, m := range s.data {
		sM = append(sM, m.ConvertToMetrics())
	}
	return sM, nil
}

func (s *MemStorage) ShowMetrics(context.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.logger.Println("Show metrics")
	lgr := logger.NewLogger("", 0)
	for _, m := range s.data {
		lgr.Println(m)
	}
}
