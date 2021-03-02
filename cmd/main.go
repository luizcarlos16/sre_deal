package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/luizcarlos16/sre_deal/cmd/get-random-number/register"

	"github.com/luizcarlos16/sre_deal/internal/config"
	"github.com/luizcarlos16/sre_deal/internal/router"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {

	config.
		Add(
			"http-service-listen-address",
			"HTTP_SERVICE_LISTEN_ADDRESS",
			string("0.0.0.0:8080"),
			"IP:PORT address to listen as service endpoint",
		).
		Add(
			"http-metrics-listen-address",
			"HTTP_METRICS_LISTEN_ADDRESS",
			string("0.0.0.0:9090"),
			"IP:PORT address to listen as metrics endpoint",
		)

}

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "myapp_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

var (
	requestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_count_total",
		Help: "Counter of HTTP requests made.",
	}, []string{"code", "method"})
	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "A histogram of latencies for requests.",
		Buckets: append([]float64{0.000001, 0.001, 0.003}, prometheus.DefBuckets...),
	}, []string{"code", "method"})
	responseSize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_size_bytes",
		Help:    "A histogram of response sizes for requests.",
		Buckets: []float64{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20},
	}, []string{"code", "method"})
)

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	config.Register().Load()

	serviceServer := http.NewServeMux()
	serviceServer.Handle("/", router.Router2)
	serviceServer.Handle("/random-number", router.Router1)

	//func que mostra o endereço no logs do docker
	go func() {
		address := config.Get("http-service-listen-address").GetStringVal()

		log.Printf("service server started on: http://%s\n", address)
		if err := http.ListenAndServe(address, serviceServer); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	metricsServer := mux.NewRouter()
	metricsServer.Use(prometheusMiddleware)
	metricsServer.Path("/metrics").Handler(promhttp.Handler())

	//func que mostra o endereço no logs do docker
	go func() {
		address := config.Get("http-metrics-listen-address").GetStringVal()

		log.Printf("metrics server started on: http://%s\n", address)
		if err := http.ListenAndServe(address, metricsServer); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//func que mostra o endereço no logs do docker
	go func() {
		sig := <-sigs
		log.Printf("received signal: %v\n", sig)
		done <- true
	}()

	log.Println("awaiting signal")
	<-done
	log.Println("exiting")
}
