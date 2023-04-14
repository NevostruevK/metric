package metrics

import (
	"errors"
	"fmt"
)

type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Параметр принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func (m *Metrics) NewGaugeMetric(id string, f float64) MetricCreater {
	return &Metrics{ID: id, MType: Gauge, Value: &f}
}

func (m *Metrics) NewCounterMetric(id string, i int64) MetricCreater {
	return &Metrics{ID: id, MType: Counter, Delta: &i}
}

func (m *Metrics) ConvertToMetrics() Metrics {
	return *m
}

func (m *Metrics) Name() string {
	return m.ID
}
func (m *Metrics) Type() string {
	return m.MType
}

func (m *Metrics) CounterValue() int64 {
	if m == nil {
		return 0
	}
	if m.Delta == nil {
		return 0
	}
	return *m.Delta
}

func (m *Metrics) GaugeValue() float64 {
	if m == nil {
		return 0
	}
	if m.Value == nil {
		return 0
	}
	return *m.Value
}

func (m *Metrics) AddCounterValue(value int64) error {
	if m.MType != Counter {
		return errors.New("error: try to add to not counter metric")
	}
	i := *m.Delta + value
	m.Delta = &i
	return nil
}

func (m Metrics) StringValue() string {
	if m.MType == Gauge {
		return roundGauge(float64(*m.Value))
	}
	return fmt.Sprintf("%d", *m.Delta)
}

func (m Metrics) String() string {
	return m.Name() + " : " + m.Type() + " : " + m.StringValue() + " : " + m.Hash
}
