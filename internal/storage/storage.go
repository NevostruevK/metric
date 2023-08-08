// Package storage организация хранилища в памяти.
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

// RepositoryData интерфейс объектов для Storage.
type RepositoryData interface {
	// Name имя объекта.
	Name() string
	// Type тип объекта.
	Type() string
	// StringValue строковое представление значения объекта.
	StringValue() string
	// CounterValue значение объекта типа "сounter".
	CounterValue() int64
	// GaugeValue значение объекта типа "gauge".
	GaugeValue() float64
	// AddCounterValue прибавление значения к объекту типа "сounter".
	AddCounterValue(int64) error
	// ConvertToMetrics преобразование к типу Metrics.
	ConvertToMetrics() metrics.Metrics
}

// Repository ннтерфейс для работы со Storage.
type Repository interface {
	// AddMetric добавление объекта RepositoryData.
	AddMetric(context.Context, RepositoryData) error
	// GetMetric чтение объекта name типа reqType.
	GetMetric(ctx context.Context, reqType, name string) (RepositoryData, error)
	// GetAllMetrics чтение всех объектов.
	GetAllMetrics(context.Context) ([]metrics.Metrics, error)
	// AddGroupOfMetrics добавление []Мetrics.
	AddGroupOfMetrics(ctx context.Context, sM []metrics.Metrics) error
	// Ping проверка коннекта к Storage.
	Ping() error
	// Close освобождение ресурсов
	Close(context.Context) error
	// SaveAllIntoFile сохранение данных в файл
	SaveAllIntoFile() (int, error)
}

// MemStorage структура для хранения метрик в памяти.
type MemStorage struct {
	Float           map[string]float64
	Int             map[string]int64
	saver           *saver
	logger          *log.Logger
	needToSyncWrite bool
	mu              sync.RWMutex
	init            bool
}

// NewMemStorage конструктор создания MemStorage.
// Параметр restore = true инициирует загрузку метрик из файла filename.
// Параметр needToSyncWrite = true инициирует синхронную запись в файл filename.
func NewMemStorage(restore, needToSyncWrite bool, filename string) *MemStorage {
	lgr := logger.NewLogger("mem storage : ", log.LstdFlags|log.Lshortfile)
	mFloat := make(map[string]float64, initialFloatMapSize)
	mInt := make(map[string]int64, initialIntMapSize)
	s, err := NewSaver(filename)
	if err != nil {
		lgr.Printf("Can't write metrics to %s\n", filename)
		return &MemStorage{Float: mFloat, Int: mInt, saver: s, logger: lgr, needToSyncWrite: false, init: true}
	}
	if filename == "" {
		return &MemStorage{Float: mFloat, Int: mInt, saver: s, logger: lgr, needToSyncWrite: false, init: true}
	}
	if restore {
		l, err := NewLoader(filename)
		if err != nil {
			lgr.Printf("Can't load metrics from %s\n", filename)
		} else {
			defer func() {
				err = l.Close()
			}()
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
	return &MemStorage{Float: mFloat, Int: mInt, saver: s, logger: lgr, needToSyncWrite: needToSyncWrite, init: true}
}

// SaveAllIntoFile сохранение всех метрик в файл.
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

func (s *MemStorage) Close(ctx context.Context) error {
	if !s.init {
		return nil
	}
	s.init = false
	count, err := s.SaveAllIntoFile()
	if err != nil {
		return err
	}
	s.logger.Printf("saved %d metrics\n", count)

	s.ShowMetrics(ctx)
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
	var errNotMetricType = errors.New("is not a metric type")

	s.mu.Lock()
	defer s.mu.Unlock()
	switch rt.Type() {
	case metrics.Counter:
		if err := rt.AddCounterValue(s.Int[rt.Name()]); err != nil {
			s.logger.Println(err)
		}
		s.Int[rt.Name()] = rt.CounterValue()
	case metrics.Gauge:
		s.Float[rt.Name()] = rt.GaugeValue()
	default:
		return errNotMetricType
	}
	if s.needToSyncWrite {
		if err := s.saver.WriteMetric(rt); err != nil {
			s.logger.Println(err)
		}
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

func PrepareMetricsForStorage(sM []metrics.Metrics) ([]metrics.Metrics, error) {
	st := NewMemStorage(false, false, "")
	if err := st.AddGroupOfMetrics(context.Background(), sM); err != nil {
		return nil, err
	}
	pM, err := st.GetAllMetrics(context.Background())
	if err != nil {
		return nil, err
	}
	return pM, nil
}
