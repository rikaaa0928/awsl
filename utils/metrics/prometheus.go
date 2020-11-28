package metrics

import (
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/global"
)

func StartMetrics(c config.Configs) {
	var err error
	global.MetricsPort, err = c.GetInt("metrics")
	if err != nil || global.MetricsPort == 0 {
		return
	}
	l, err := net.Listen("tcp", ":"+strconv.FormatInt(global.MetricsPort, 10))
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	err = http.Serve(l, mux)
	if err != nil {
		log.Println(err)
	}
}
