package metrics

import (
	"errors"
	"fmt"
	"strconv"
)

type Metric struct {
	MName   string	
	MType  string	
	GValue gauge	
	CValue counter	
}

func NewGaugeMetric(name string, f float64) *Metric {
	return &Metric{MName: name, MType: Gauge, GValue: gauge(f)}
}

func NewCounterMetric(name string, i int64) *Metric {
	return &Metric{MName: name, MType: Counter, CValue: counter(i)}
}

func (m *Metric) NewGaugeMetric(name string, value float64)  MetricCreater{
	return NewGaugeMetric(name, value)
}

func (m *Metric) NewCounterMetric(name string, value int64)  MetricCreater{
	return NewCounterMetric(name, value)
}

func NewValueMetric(name string, typeM string, value string) (*Metric, error) {
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

func (m *Metric) Name() string {
	return m.MName
}

func (m *Metric) Type() string {
	return m.MType
}

func (m *Metric) CounterValue() int64 {
	return int64(m.CValue)
}

func (m *Metric) AddCounterValue(value int64)  error {
	if m.MType != Counter {
		return errors.New("error: try to add to not counter metric")
	}
	m.CValue += counter(value)
	return nil
}


func (m Metric) StringValue() string {
	if m.MType == Gauge {
		return fmt.Sprintf("%.3f", float64(m.GValue))
	}
	return fmt.Sprintf("%d", m.CValue)
}

func (m Metric) String() string {
	return m.Type() + "/" + m.Name() + "/" + m.StringValue()
}
