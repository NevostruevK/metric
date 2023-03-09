package metrics

import (
	"errors"
	"strconv"
	"strings"
)

func URLToMetric(url string) (*Metric, error) {
	words := strings.Split(url, "/")
/*		for idx, word := range words {
			fmt.Printf("Word %d is: %s\n", idx, word)
		}
*/	if len(words) != 5 {
		return nil, errors.New("wrong slash count error")
	}
	if words[0] != "" && words[1] != "update" {
		return nil, errors.New("wrong prefix error")
	}

	switch words[2] {
	case "gauge":
		f, err := strconv.ParseFloat(words[4], 64)
		if err != nil {
			return nil, errors.New("parse to gauge error")
		}
		return NewGaugeMetric(words[3], f), nil

	case "counter":
		i, err := strconv.ParseInt(words[4], 10, 64)
		if err != nil {
			return nil, errors.New("parse to counter error")
		}
		return NewCounterMetric(words[3], i), nil
	default:
		return nil, errors.New("type error")
	}
}