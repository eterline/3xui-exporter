package metrics

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/constraints"
)

func joinParams(v ...any) string {
	vStr := make([]string, len(v))
	for i, value := range v {
		vStr[i] = fmt.Sprintf("%v", value)
	}
	return strings.Join(vStr, "/")
}

func setMetric[T constraints.Integer | constraints.Float](p *prometheus.GaugeVec, value T, params ...string) {
	p.WithLabelValues(params...).Set(convertNumberToFloat(value))
}

func convertNumberToFloat[T constraints.Integer | constraints.Float](v T) float64 {
	switch any(v).(type) {
	case int:
		return float64(any(v).(int))
	case int8:
		return float64(any(v).(int8))
	case int16:
		return float64(any(v).(int16))
	case int32:
		return float64(any(v).(int32))
	case int64:
		return float64(any(v).(int64))
	case uint:
		return float64(any(v).(uint))
	case uint8:
		return float64(any(v).(uint8))
	case uint16:
		return float64(any(v).(uint16))
	case uint32:
		return float64(any(v).(uint32))
	case uint64:
		return float64(any(v).(uint64))
	case float32:
		return float64(any(v).(float32))
	case float64:
		return any(v).(float64)
	default:
		return 0.0
	}
}

func newGaugeVec(name, help string, tags []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		tags,
	)
}

func newCounterVec(name, help string, tags []string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		tags,
	)
}

func strValueIs(v string, eq ...string) float64 {
	for _, val := range eq {
		if v == val {
			return 1.0
		}
	}
	return 0.0
}
