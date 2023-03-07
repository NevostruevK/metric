package metrics

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const MetricsCount = 29

type gauge float64
type counter int64

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
        NewMetric(s string) Metric
}
func (m Metric) Name() string {
        return m.name
}
func (m Metric) String() string {
        s := m.typeM + "/" + m.name + "/"
        if m.typeM == "gauge" {
                //              return s + fmt.Sprintf("%.2f",float64(m.gValue))
                return s + fmt.Sprintf("%f", float64(m.gValue))
        }
        return s + fmt.Sprintf("%d", m.cValue)
}

func URLToMetric(url string) (*Metric, error){
	words := strings.Split(url, "/")
	for idx, word := range words {
		fmt.Printf("Word %d is: %s\n", idx, word)
	}
        if len(words) != 5{
                return nil, errors.New("wrong error")
        }
        if words[0] != "" && words[1] != "update"{
                return nil, errors.New("wrong prefix")
        }

        switch words[2]{
        case "gauge":
                f, err := strconv.ParseFloat(words[4], 64) 
                if err != nil{
                        return nil, errors.New("parse to gauge error")
                }
//                return NewGaugeMetric("123", 0.123), nil
                return NewGaugeMetric(words[3], f), nil
        
        case "counter":
                i, err := strconv.ParseInt(words[4], 10, 64) 
                if err != nil{
                        return nil, errors.New("parse to counter error")
                }
                return NewCounterMetric(words[3], i), nil
        default:
                return nil, errors.New("type error")
        }
}