package metrics

import (
	"errors"
	"fmt"
	"strconv"
)

type Metric struct {
	name   string
	typeM  string
	gValue gauge
	cValue counter
}

func NewGaugeMetric(name string, f float64) *Metric {
	return &Metric{name: name, typeM: Gauge, gValue: gauge(f)}
}

func NewCounterMetric(name string, i int64) *Metric {
	return &Metric{name: name, typeM: Counter, cValue: counter(i)}
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

func (m Metric) Name() string {
	return m.name
}

func (m Metric) Type() string {
	return m.typeM
}


func (m Metric) StringValue() string {
	if m.typeM == Gauge {
		return fmt.Sprintf("%.3f", float64(m.gValue))
	}
	return fmt.Sprintf("%d", m.cValue)
}

func (m Metric) String() string {
	return m.Type() + "/" + m.Name() + "/" + m.StringValue()
}

func (m *Metric) AddMetricValue(new Metric) (*Metric, error) {
	if m.typeM != new.typeM {
		return m, errors.New("error: try to add different types")
	}
	m.cValue += new.cValue
	m.gValue += new.gValue
	return m, nil
}