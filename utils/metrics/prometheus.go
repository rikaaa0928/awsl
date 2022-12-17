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
	log.Println(global.MetricsPort)
	if err != nil || global.MetricsPort == 0 {
		return
	}
	l, err := net.Listen("tcp", ":"+strconv.FormatInt(global.MetricsPort, 10))
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	//mux.HandleFunc("/debug/pprof/", pprof.Index)
	//mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	//mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	//mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	//mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	err = http.Serve(l, mux)
	if err != nil {
		log.Println(err)
	}
}
