package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "", // metricas de quantas requições q a app teve
			Help: "Number of requests",
		}, []string{"path"}) //dimensão (endpoint)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of http request",
		}, []string{"path"})
)

func init() {
	//registro de métricas no prometheus
	prometheus.MustRegister(requestDuration)
}
func handler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.URL.Path))
	defer timer.ObserveDuration()
	httpRequests.WithLabelValues(r.URL.Path)

	w.Write([]byte("OK"))
}
func main() {
	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Listening on :8080")

	http.ListenAndServe(":8080", nil)
}
