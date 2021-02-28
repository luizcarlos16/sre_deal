package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/luizcarlos16/sre_deal/cmd/get-random-number/register"

	"github.com/luizcarlos16/sre_deal/internal/config"
	"github.com/luizcarlos16/sre_deal/internal/router"
	
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	
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

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(responseSize)
}

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	fmt.Fprintf(w, "OK\n")
}

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
	
	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	wrapHandler := promhttp.InstrumentHandlerCounter(
		requestCount,
		promhttp.InstrumentHandlerDuration(
			requestDuration,
			promhttp.InstrumentHandlerResponseSize(responseSize, http.HandlerFunc(handler)),
		),
	)
	http.Handle("/", wrapHandler)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("metrics server started on: https://localhost:9090")

	http.ListenAndServe(":9090", nil)
	

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