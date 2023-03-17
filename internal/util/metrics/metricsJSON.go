package metrics

import "fmt"

type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Параметр принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}

func (m *Metrics) NewGaugeMetric(id string, f float64) MetricCreater{
	return &Metrics{ID: id, MType: Gauge, Value: &f}
}

func (m Metrics) NewCounterMetric(id string, i int64) MetricCreater{
	return &Metrics{ID: id, MType: Counter, Delta: &i}
}

func (m Metrics) Name() string {
	return m.ID
}
func (m Metrics) Type() string {
	return m.MType
}

func (m Metrics) StringValue() string {
	if m.MType == Gauge {
		return fmt.Sprintf("%v", m.Value)
	}
	return fmt.Sprintf("%v", m.Delta)
}