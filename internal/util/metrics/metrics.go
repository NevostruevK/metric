package metrics

import (
	"fmt"
)

type gauge float64
type counter int64
func (g gauge) NewMetric(name string) Metric{
	return Metric{name: name, typeM: GAUGE, gValue: g}
}
func (c counter) NewMetric(name string) Metric{
	return Metric{name: name, typeM: COUNTER, cValue: c}
}
func Float64ToGauge(f float64) gauge{
	return gauge(f)
}
func Int64ToGauge(d int64) counter{
	return counter(d)
}

type MetricType string
type MetricName string
const (
	GAUGE  = string("gauge")
	COUNTER = string("counter")
)
type Metric struct{
	name string
	typeM string
	gValue gauge
	cValue counter
}
type NewMetricItn interface{
	NewMetric(s string) Metric
}

func (m Metric) String() string{
	s := m.typeM + "/" + m.name + "/"
	if m.typeM == GAUGE{
//		return s + fmt.Sprintf("%.2f",float64(m.gValue))
		return s + fmt.Sprintf("%f",float64(m.gValue))
	}	
	return s + fmt.Sprintf("%d",m.cValue)
}
/*
func (m Metric) GetType() string{
	return string(m.typeM)
}
func (m Metric) GetName() string{
	return string(m.name)
}
*/
func (m Metric) ValueToString() string{
	if (m.typeM==COUNTER){
		return fmt.Sprintf("%d", m.cValue)
	}
//	return fmt.Sprintf("%.2f", m.gValue)
	return fmt.Sprintf("%f", m.gValue)
//	string(m.GetValue())
}
const MetricsCount = 29