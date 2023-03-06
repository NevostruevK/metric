package metrics

import (
	"fmt"
)

const MetricsCount = 29

type gauge float64
type counter int64

func (g gauge) NewMetric(name string) Metric {
        return Metric{name: name, typeM: "gauge", gValue: g}
}
func (c counter) NewMetric(name string) Metric {
        return Metric{name: name, typeM: "counter", cValue: c}
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
        NewMetric(s string) Metric
}

func (m Metric) String() string {
        s := m.typeM + "/" + m.name + "/"
        if m.typeM == "gauge" {
                //              return s + fmt.Sprintf("%.2f",float64(m.gValue))
                return s + fmt.Sprintf("%f", float64(m.gValue))
        }
        return s + fmt.Sprintf("%d", m.cValue)
}