package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "yorkie"
	subsystem = "rpcserver"
)

type RPCServerMetrics struct {
	namespace string

	pushpullResponseSeconds prometheus.Histogram
}

func NewRPCServerMetrics() *RPCServerMetrics {
	metrics := &RPCServerMetrics{
		namespace: namespace,
	}
	metrics.recordMetrics()
	return metrics
}

func (r *RPCServerMetrics) recordMetrics() {
	r.pushpullResponseSeconds = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "pushpull_response_seconds",
		Help:      "Response time of PushPull API.",
	})
}

func (r *RPCServerMetrics) ObservePushpullResponseSeconds(seconds float64) {
	r.pushpullResponseSeconds.Observe(seconds)
}

func (r *RPCServerMetrics) AddPushpullReceivedChanges(count float64) {

}

func (r *RPCServerMetrics) AddPushpullSentChanges(count float64) {

}

func (r *RPCServerMetrics) ObservePushpullSnapshotDurationSeconds(seconds float64) {

}

func (r *RPCServerMetrics) AddPushpullSnapshotBytes(byte float64) {

}
