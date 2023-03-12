package metrics

import (
	"errors"
	"fmt"
	"strconv"
)

const MetricsCount = 29

type gauge float64
type counter int64

const (
	Gauge   = "gauge"
	Counter = "counter"
)

func IsMetricType(checkType string) bool {
	if checkType != Gauge && checkType != Counter {
		return false
	}
	return true
}

type Metric struct {
	name   string
	typeM  string
	gValue gauge
	cValue counter
}

type NewMetricItn interface {
	NewMetric(name string) Metric
}

func (g gauge) NewMetric(name string) *Metric {
	return &Metric{name: name, typeM: "gauge", gValue: g}
}
func (c counter) NewMetric(name string) *Metric {
	return &Metric{name: name, typeM: "counter", cValue: c}
}
func NewGaugeMetric(name string, f float64) *Metric {
	return gauge(f).NewMetric(name)
}
func NewCounterMetric(name string, i int64) *Metric {
	return counter(i).NewMetric(name)
}
func NewValueMetric(name string, typeM string, value string) (*Metric, error) {
	switch typeM {
	case "gauge":
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, errors.New("convert to gauge with an error")
		}
		return NewGaugeMetric(name, f), nil
	case "counter":
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

func (m Metric) Value() string {
	if m.typeM == "gauge" {
		return fmt.Sprintf("%.3f", float64(m.gValue))
	}
	return fmt.Sprintf("%d", m.cValue)
}

func (m Metric) String() string {
	return m.Type() + "/" + m.Name() + "/" + m.Value()
}

func (m Metric) AddMetricValue(new Metric) (Metric, error) {
	if m.typeM != new.typeM {
		return m, errors.New("error: try to add differnt types")
	}
	m.cValue += new.cValue
	m.gValue += new.gValue
	return m, nil
}
