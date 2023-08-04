package metrics

import (
	"errors"
	"fmt"

	pb "github.com/NevostruevK/metric/proto"
)

// Metrics структура метрики для конвертации в JSON.
type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Параметр принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

// NewJSONGaugeMetric конструктор метрики id со значение f типа "gauge".
func NewJSONGaugeMetric(id string, f float64) Metrics {
	return Metrics{ID: id, MType: Gauge, Value: &f}
}

// NewJSONCounterMetric конструктор метрики id со значение i типа "counter".
func NewJSONCounterMetric(id string, i int64) Metrics {
	return Metrics{ID: id, MType: Counter, Delta: &i}
}

// ToProto конвертит метрику для gRPC.
func (m *Metrics) ToProto() *pb.Metric {
	var pbMetric pb.Metric
	if m.MType == Gauge {
		pbMetric.Type = pb.MetricType_GAUGE
		pbMetric.Value = m.GaugeValue()
	} else {
		pbMetric.Type = pb.MetricType_COUNTER
		pbMetric.Delta = m.CounterValue()
	}
	pbMetric.Name = m.ID
	pbMetric.Hash = m.Hash
	return &pbMetric
}

// ToProto конвертирует массив метрик для gRPC.
func ToProto(sm []Metrics) []*pb.Metric {
	sp := make([]*pb.Metric, len(sm))
	for i, m := range sm {
		sp[i] = m.ToProto()
	}
	return sp
}

// NewGaugeMetric метод создания метрики id со значение f типа "gauge".
func (m *Metrics) NewGaugeMetric(id string, f float64) Metrics {
	m.ID = id
	m.Value = &f
	m.MType = Gauge
	return *m
}

// NewCounterMetric метод создания метрики id со значение i типа "counter".
func (m *Metrics) NewCounterMetric(id string, i int64) Metrics {
	m.ID = id
	m.Delta = &i
	m.MType = Counter
	return *m
}

// ConvertToMetrics заглушка для удовлетворения интерфейсу RepositoryData.
func (m *Metrics) ConvertToMetrics() Metrics {
	return *m
}

// Name возвращает имя метрики.
func (m *Metrics) Name() string {
	return m.ID
}

// Type возвращает тип метрики.
func (m *Metrics) Type() string {
	return m.MType
}

// CounterValue возвращает значение метрики типа "сounter".
func (m *Metrics) CounterValue() int64 {
	if m == nil {
		return 0
	}
	if m.Delta == nil {
		return 0
	}
	return *m.Delta
}

// GaugeValue возвращает значение метрики типа "gauge".
func (m *Metrics) GaugeValue() float64 {
	if m == nil {
		return 0
	}
	if m.Value == nil {
		return 0
	}
	return *m.Value
}

// AddCounterValue прибавление значения к метрике типа "сounter".
func (m *Metrics) AddCounterValue(value int64) error {
	if m.MType != Counter {
		return errors.New("error: try to add to not counter metric")
	}
	i := *m.Delta + value
	m.Delta = &i
	return nil
}

// StringValue строковое представление значения метрики.
func (m Metrics) StringValue() string {
	if m.MType == Gauge {
		return roundGauge(float64(*m.Value))
	}
	return fmt.Sprintf("%d", *m.Delta)
}

// String строковое представление метрики.
func (m Metrics) String() string {
	return m.Type() + " : " + m.Name() + " : " + m.StringValue() + " : " + m.Hash
}

// StringWithSlash строковое представление метрики, разделенное "/".
func (m Metrics) StringWithSlash() string {
	return m.Type() + "/" + m.Name() + "/" + m.StringValue()
}
