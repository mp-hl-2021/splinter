package v1

import (
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type Api struct {
	useCases usecases.UserInterface
}

func NewApi(u usecases.UserInterface) *Api {
	return &Api{useCases: u}
}

func (a *Api) Router(router *mux.Router) {
	router.HandleFunc("/create_account", a.endpointCreateAccount).Methods(http.MethodPost)
	router.HandleFunc("/authenticate", a.endpointAuthenticate).Methods(http.MethodPost)
	router.HandleFunc("/users/current", a.endpointGetCurrentUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{user}", a.endpointGetUser).Methods(http.MethodGet)

	router.HandleFunc("/snippets", a.endpointPostSnippet).Methods(http.MethodPost)
	router.HandleFunc("/users/{user}/snippets", a.endpointGetSnippetsByUser).Methods(http.MethodGet)
	router.HandleFunc("/snippets/language/{language}", a.endpointGetSnippetsByLanguage).Methods(http.MethodGet)
	router.HandleFunc("/snippets/{snippet}", a.endpointGetSnippet).Methods(http.MethodGet)
	router.HandleFunc("/snippets/{snippet}", a.endpointDeleteSnippet).Methods(http.MethodDelete)
	router.HandleFunc("/snippets/{snippet}/vote", a.endpointDeleteSnippet).Methods(http.MethodDelete)

	router.HandleFunc("/snippets/{snippet}/comments", a.endpointGetComments).Methods(http.MethodGet)
	router.HandleFunc("/snippets/{snippet}/comments", a.endpointPostComment).Methods(http.MethodPost)
	router.HandleFunc("/comments/{comment}", a.endpointDeleteComment).Methods(http.MethodDelete)
}
