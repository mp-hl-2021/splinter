package v1

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/auth"
	"github.com/mp-hl-2021/splinter/types"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type Api struct {
	useCases      usecases.UserInterface
	authenticator auth.Authenticator
}

func NewApi(u usecases.UserInterface, a auth.Authenticator) *Api {
	return &Api{useCases: u, authenticator: a}
}

func makeAuthMiddleware(a auth.Authenticator) func (handler http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			uid, err := a.UserIdByToken(token)
			if err != nil {
				WriteError(w, err, http.StatusForbidden)
				return
			}
			context.Set(r, "uid", types.UserId(uid))
			next.ServeHTTP(w, r)
		})
	}
}

func (a *Api) Router(router *mux.Router) {
	amw := makeAuthMiddleware(a.authenticator)
	router.HandleFunc("/create_account", a.endpointCreateAccount).Methods(http.MethodPost)
	router.HandleFunc("/authenticate", a.endpointAuthenticate).Methods(http.MethodPost)

	router.Handle("/users/current", amw(http.HandlerFunc(a.endpointGetCurrentUser))).Methods(http.MethodGet)
	router.Handle("/users/{user}", amw(http.HandlerFunc(a.endpointGetUser))).Methods(http.MethodGet)

	router.Handle("/snippets", amw(http.HandlerFunc(a.endpointPostSnippet))).Methods(http.MethodPost)
	router.Handle("/users/{user}/snippets", amw(http.HandlerFunc(a.endpointGetSnippetsByUser))).Methods(http.MethodGet)
	router.Handle("/snippets/language/{language}", amw(http.HandlerFunc(a.endpointGetSnippetsByLanguage))).Methods(http.MethodGet)
	router.Handle("/snippets/{snippet}", amw(http.HandlerFunc(a.endpointGetSnippet))).Methods(http.MethodGet)
	router.Handle("/snippets/{snippet}", amw(http.HandlerFunc(a.endpointDeleteSnippet))).Methods(http.MethodDelete)
	router.Handle("/snippets/{snippet}/vote", amw(http.HandlerFunc(a.endpointVote))).Methods(http.MethodPost)

	router.Handle("/snippets/{snippet}/comments", amw(http.HandlerFunc(a.endpointGetComments))).Methods(http.MethodGet)
	router.Handle("/snippets/{snippet}/comments", amw(http.HandlerFunc(a.endpointPostComment))).Methods(http.MethodPost)
	router.Handle("/comments/{comment}", amw(http.HandlerFunc(a.endpointDeleteComment))).Methods(http.MethodDelete)
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

func GetCurrentUid(r *http.Request) types.UserId {
	return context.Get(r, "uid").(types.UserId)
}