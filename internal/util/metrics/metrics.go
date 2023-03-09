package metrics

import (
	"errors"
	"fmt"
)

const MetricsCount = 29

type gauge float64
type counter int64

const(
        Gauge = "gauge"
        Counter = "counter"
)

func (g gauge) NewMetric(name string) *Metric {
        return &Metric{name: name, typeM: "gauge", gValue: g}
}
func (c counter) NewMetric(name string) *Metric {
        return &Metric{name: name, typeM: "counter", cValue: c}
}
func NewGaugeMetric(name string, f float64) *Metric{
        return gauge(f).NewMetric(name)
}
func NewCounterMetric(name string, i int64) *Metric{
        return counter(i).NewMetric(name)
}

func Float64ToGauge(f float64) gauge {
        return gauge(f)
}
func Int64ToGauge(d int64) counter {
        return counter(d)
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
func (m Metric) Name() string {
        return m.name
}

func (m Metric) Type() string {
        return m.typeM
}

func (m Metric) String() string {
        s := m.typeM + "/" + m.name + "/"
        if m.typeM == "gauge" {
                //              return s + fmt.Sprintf("%.2f",float64(m.gValue))
                return s + fmt.Sprintf("%f", float64(m.gValue))
        }
        return s + fmt.Sprintf("%d", m.cValue)
}

func (m Metric) AddMetricValue(new Metric) (Metric, error){
        if m.typeM != new.typeM{
                return m, errors.New("error try to add differnt types")
        }
        m.cValue += new.cValue
        m.gValue += new.gValue
        return m, nil
}