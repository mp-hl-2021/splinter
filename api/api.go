package api

import (
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/api/v1"
	"github.com/mp-hl-2021/splinter/auth"
	"github.com/mp-hl-2021/splinter/monitoring"
	"github.com/mp-hl-2021/splinter/usecases"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Api struct {
	v1 *v1.Api
}

func NewApi(u usecases.UserInterface, a auth.Authenticator) *Api {
	return &Api{v1: v1.NewApi(u, a)}
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	a.v1.Router(router.PathPrefix("/api/v1").Subrouter())
	router.Handle("/metrics", promhttp.Handler())
	router.Use(monitoring.Measurer())

	return router
}
