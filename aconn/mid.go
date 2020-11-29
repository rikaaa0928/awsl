package aconn

import (
	"context"

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
	}, []string{"type", "tag", "end_addr"})
	prometheus.MustRegister(realTimeConnectionNum)
}

func NewMetricsMid(ctx context.Context, tag, typ, endAddr string) MetricsMid {
	if global.MetricsPort > 0 {
		realTimeConnectionNum.With(prometheus.Labels{"type": typ, "tag": tag, "end_addr": endAddr}).Inc()
	}
	return MetricsMid{typ: typ, tag: tag, endAddr: endAddr}
}

func NewMetricsMidWithCTX(ctx context.Context, conf map[string]interface{}, endAddr string) MetricsMid {
	if global.MetricsPort > 0 {
		outTag := ctx.Value(global.CTXRoute)
		if outTag == nil {
			return MetricsMid{disabled: true}
		}
		tag, ok := outTag.(string)
		if !ok {
			return MetricsMid{disabled: true}
		}
		tagConf := conf[tag].(map[string]interface{})
		typ := tagConf["type"].(string)
		realTimeConnectionNum.With(prometheus.Labels{"type": typ, "tag": tag, "end_addr": endAddr}).Inc()
		return MetricsMid{typ: typ, tag: tag, endAddr: endAddr}
	}
	return MetricsMid{disabled: true}
}

type MetricsMid struct {
	typ      string
	tag      string
	endAddr  string
	disabled bool
}

func (m MetricsMid) MetricsClose(next Closer) Closer {
	return func() error {
		if global.MetricsPort > 0 && !m.disabled {
			realTimeConnectionNum.With(prometheus.Labels{"type": m.typ, "tag": m.tag, "end_addr": m.endAddr}).Dec()
		}
		return next()
	}
}
