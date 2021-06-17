package api

import (
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/api/v1"
	"github.com/mp-hl-2021/splinter/auth"
	"github.com/mp-hl-2021/splinter/monitoring"
	"github.com/mp-hl-2021/splinter/usecases"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type Api struct {
	v1 *v1.Api
}

func NewApi(u usecases.UserInterface, a auth.Authenticator) *Api {
	return &Api{v1: v1.NewApi(u, a)}
}

type responseWriterObserver struct {
	http.ResponseWriter
	status int
}

func (o *responseWriterObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	o.status = code
}

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		ow := &responseWriterObserver{w, http.StatusOK}
		next.ServeHTTP(ow, r)
		log.Printf("method = %v, url = %v, status = %v", r.Method, r.URL.Path, ow.status)
	})
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	router.Use(logRequestMiddleware)

	a.v1.Router(router.PathPrefix("/api/v1").Subrouter())
	router.Handle("/metrics", promhttp.Handler())
	router.Use(monitoring.Measurer())

	return router
}
