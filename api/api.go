package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/api/v1"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type Api struct {
	v1 *v1.Api
}

func NewApi(u usecases.UserInterface) *Api {
	return &Api{v1: v1.NewApi(u)}
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	a.v1.Router(router.PathPrefix("/api/v1").Subrouter())

	return router
}

type errorResponse struct {
	StatusCode int
	Error      string
}

func WriteError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{
		StatusCode: statusCode,
		Error:      err.Error(),
	})
}
