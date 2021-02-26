package router
  
import (
  "net/http"

  "github.com/gorilla/mux"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promauto"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
  httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
    Name: "myapp_http_duration_seconds",
    Help: "Duration of HTTP requests.",
  }, []string{"path"})
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    Router := mux.CurrentRoute(r)
    path, _ := Router.GetPathTemplate()
    timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
    next.ServeHTTP(w, r)
    timer.ObserveDuration()
  })
}

func main() {
  Router := mux.NewRouter()
  Router.Use(prometheusMiddleware)
  Router.Path("/metrics").Handler(promhttp.Handler())
  Router.Path("/obj/{id}").HandlerFunc(
    func(w http.ResponseWriter, Router *http.Request) {})

  srv := &http.Server{Addr: "localhost:1234", Handler: Router}
  srv.ListenAndServe()
}