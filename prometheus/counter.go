package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
)

type CounterVec struct {
	cv *prometheus.CounterVec
}

func NewCounterVec(opts prometheus.CounterOpts, keys []attribute.Key) *CounterVec {
	labelNames := make([]string, 0, len(keys))
	for _, k := range keys {
		labelNames = append(labelNames, string(k))
	}
	return &CounterVec{
		cv: prometheus.NewCounterVec(opts, labelNames),
	}
}

func (v *CounterVec) With(labels ...attribute.KeyValue) prometheus.Counter {
	return v.cv.With(attributesToLabels(labels...))
}

func (v *CounterVec) Describe(ch chan<- *prometheus.Desc) {
	v.cv.Describe(ch)
}

func (v *CounterVec) Collect(ch chan<- prometheus.Metric) {
	v.cv.Collect(ch)
}

func attributesToLabels(attrs ...attribute.KeyValue) prometheus.Labels {
	ls := make(map[string]string, len(attrs))
	for _, a := range attrs {
		ls[string(a.Key)] = a.Value.Emit()
	}
	return ls
}
