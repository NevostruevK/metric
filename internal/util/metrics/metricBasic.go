package metrics

import (
	"errors"
	"fmt"
	"strconv"
)

// BasicMetric структура метрики.
type BasicMetric struct {
	MName  string
	MType  string
	GValue gauge
	CValue counter
}

// NewGaugeMetric конструктор BasicMetric name со значение f типа "gauge".
func NewGaugeMetric(name string, f float64) *BasicMetric {
	return &BasicMetric{MName: name, MType: Gauge, GValue: gauge(f)}
}

// NewCounterMetric конструктор BasicMetric name со значение i типа "counter".
func NewCounterMetric(name string, i int64) *BasicMetric {
	return &BasicMetric{MName: name, MType: Counter, CValue: counter(i)}
}

// NewGaugeMetric метод создания метрики name со значение value типа "gauge".
func (m *BasicMetric) NewGaugeMetric(name string, value float64) MetricCreater {
	return NewGaugeMetric(name, value)
}

// NewCounterMetric метод создания метрики name со значение value типа "counter".
func (m *BasicMetric) NewCounterMetric(name string, value int64) MetricCreater {
	return NewCounterMetric(name, value)
}

// NewValueMetric конструктор BasicMetric name со значение value типа typeM.
func NewValueMetric(name string, typeM string, value string) (*BasicMetric, error) {
	switch typeM {
	case Gauge:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, errors.New("convert to gauge with an error")
		}
		return NewGaugeMetric(name, f), nil
	case Counter:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, errors.New("convert to counter with an error")
		}
		return NewCounterMetric(name, i), nil
	default:
		return nil, errors.New("type error")
	}
}

// ConvertToMetrics преобразование к типу Metrics.
func (m *BasicMetric) ConvertToMetrics() Metrics {
	if m.MType == Gauge {
		return NewJSONGaugeMetric(m.MName, float64(m.GValue))
	}
	return NewJSONCounterMetric(m.MName, int64(m.CValue))
}

// Name возвращает имя метрики.
func (m *BasicMetric) Name() string {
	return m.MName
}

// Type возвращает тип метрики.
func (m *BasicMetric) Type() string {
	return m.MType
}

// CounterValue возвращает значение метрики типа "сounter".
func (m *BasicMetric) CounterValue() int64 {
	return int64(m.CValue)
}

// GaugeValue возвращает значение метрики типа "gauge".
func (m *BasicMetric) GaugeValue() float64 {
	return float64(m.GValue)
}

// AddCounterValue прибавление значения к метрике типа "сounter".
func (m *BasicMetric) AddCounterValue(value int64) error {
	if m.MType != Counter {
		return errors.New("error: try to add to not counter metric")
	}
	m.CValue += counter(value)
	return nil
}

// StringValue строковое представление значения метрики.
func (m BasicMetric) StringValue() string {
	if m.MType == Gauge {
		return roundGauge(float64(m.GValue))
	}
	return fmt.Sprintf("%d", m.CValue)
}

// String строковое представление метрики.
func (m BasicMetric) String() string {
	return m.Type() + " " + m.Name() + " " + m.StringValue()
}

// StringWithSlash строковое представление метрики, разделенное "/".
func (m BasicMetric) StringWithSlash() string {
	return m.Type() + "/" + m.Name() + "/" + m.StringValue()
}
