package aconn

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/rikaaa0928/awsl/global"
)

var realTimeConnectionNum *prometheus.GaugeVec

func init() {
	realTimeConnectionNum = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "awsl",
		Subsystem: "aconn",
		Name:      "realtime_connection_num",
		Help:      "Number of realtime connection.",
	}, []string{"type"})
	prometheus.MustRegister(realTimeConnectionNum)
}

type MetricsMid struct {
	Typ string
}

func (m MetricsMid) MetricsClose(closer Closer) Closer {
	return func() error {
		if global.MetricsPort > 0 {
			realTimeConnectionNum.With(prometheus.Labels{"type": m.Typ}).Dec()
		}
		return closer()
	}
}
