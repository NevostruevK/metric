package metrics

import (
	"errors"
	"fmt"
	"strconv"
)

type BasicMetric struct {
	MName  string
	MType  string
	GValue gauge
	CValue counter
}

func NewGaugeMetric(name string, f float64) *BasicMetric {
	return &BasicMetric{MName: name, MType: Gauge, GValue: gauge(f)}
}

func NewCounterMetric(name string, i int64) *BasicMetric {
	return &BasicMetric{MName: name, MType: Counter, CValue: counter(i)}
}

func (m *BasicMetric) NewGaugeMetric(name string, value float64) MetricCreater {
	return NewGaugeMetric(name, value)
}

func (m *BasicMetric) NewCounterMetric(name string, value int64) MetricCreater {
	return NewCounterMetric(name, value)
}

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

func (m *BasicMetric) Name() string {
	return m.MName
}

func (m *BasicMetric) Type() string {
	return m.MType
}

func (m *BasicMetric) CounterValue() int64 {
	return int64(m.CValue)
}

func (m *BasicMetric) AddCounterValue(value int64) error {
	if m.MType != Counter {
		return errors.New("error: try to add to not counter metric")
	}
	m.CValue += counter(value)
	return nil
}

func (m BasicMetric) StringValue() string {
	if m.MType == Gauge {
		return roundGauge(float64(m.GValue))
	}
	return fmt.Sprintf("%d", m.CValue)
}

func (m BasicMetric) String() string {
	return m.Type() + "/" + m.Name() + "/" + m.StringValue()
}
