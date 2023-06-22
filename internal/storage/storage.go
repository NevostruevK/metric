package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
)

const initialFloatMapSize = 32
const initialIntMapSize = 4

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
	Float           map[string]float64
	Int             map[string]int64
	saver           *saver
	logger          *log.Logger
	needToSyncWrite bool
	mu              sync.RWMutex
}

func NewMemStorage(restore, needToSyncWrite bool, filename string) *MemStorage {
	lgr := logger.NewLogger("mem storage : ", log.LstdFlags|log.Lshortfile)
	mFloat := make(map[string]float64, initialFloatMapSize)
	mInt := make(map[string]int64, initialIntMapSize)
	s, err := NewSaver(filename)
	if err != nil {
		lgr.Printf("Can't write metrics to %s\n", filename)
		return &MemStorage{Float: mFloat, Int: mInt, saver: s, logger: lgr, needToSyncWrite: false}
	}
	if filename == "" {
		return &MemStorage{Float: mFloat, Int: mInt, saver: s, logger: lgr, needToSyncWrite: false}
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
				switch m.MType {
				case metrics.Counter:
					mInt[m.Name()] = m.CounterValue()
				case metrics.Gauge:
					mFloat[m.Name()] = m.GaugeValue()
				}
			}
		}
	}
	return &MemStorage{Float: mFloat, Int: mInt, saver: s, logger: lgr, needToSyncWrite: needToSyncWrite}
}
func (s *MemStorage) SaveAllIntoFile() (int, error) {
	if s.saver == nil {
		s.logger.Println("can't save metrics into file, saver wasn't initialized")
		return 0, fmt.Errorf("can't save metrics into file, saver wasn't initialized")
	}
	count := 0
	for name, f := range s.Float {
		if err := s.saver.WriteMetric(metrics.NewGaugeMetric(name, f)); err != nil {
			msg := fmt.Sprintf("ERROR : can't save metric into file, encoder error %v\n", err)
			s.logger.Println(msg)
			return count, fmt.Errorf(msg)
		}
		count++
	}
	for name, i := range s.Int {
		if err := s.saver.WriteMetric(metrics.NewCounterMetric(name, i)); err != nil {
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

var errNotMetricType = errors.New("is not a metric type")

func (s *MemStorage) AddMetric(ctx context.Context, rt RepositoryData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch rt.Type() {
	case metrics.Counter:
		rt.AddCounterValue(s.Int[rt.Name()])
		s.Int[rt.Name()] = rt.CounterValue()
	case metrics.Gauge:
		s.Float[rt.Name()] = rt.GaugeValue()
	default:
		return errNotMetricType
	}
	if s.needToSyncWrite {
		s.saver.WriteMetric(rt)
	}
	return nil
}

func (s *MemStorage) GetMetric(ctx context.Context, reqType, name string) (RepositoryData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	switch reqType {
	case metrics.Counter:
		i, ok := s.Int[name]
		if ok {
			return metrics.NewCounterMetric(name, i), nil
		}
	case metrics.Gauge:
		f, ok := s.Float[name]
		if ok {
			return metrics.NewGaugeMetric(name, f), nil
		}
	}
	return nil, fmt.Errorf("type %s : name %s is not valid metric type", reqType, name)
}

func (s *MemStorage) GetAllMetrics(context.Context) ([]metrics.Metrics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sM := make([]metrics.Metrics, 0, len(s.Float)+len(s.Int))
	for name, f := range s.Float {
		sM = append(sM, metrics.NewJSONGaugeMetric(name, f))
	}
	for name, i := range s.Int {
		sM = append(sM, metrics.NewJSONCounterMetric(name, i))
	}
	return sM, nil
}

func (s *MemStorage) ShowMetrics(context.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.logger.Println("Show metrics")
	lgr := logger.NewLogger("", 0)
	for name, f := range s.Float {
		lgr.Println(metrics.NewJSONGaugeMetric(name, f))
	}
	for name, i := range s.Int {
		lgr.Println(metrics.NewJSONCounterMetric(name, i))
	}
}
