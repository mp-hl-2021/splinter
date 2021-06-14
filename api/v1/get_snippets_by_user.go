package v1

// Endpoint: /api/v1/users/{user}/snippets
// Method: GET

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
	"strconv"
)

type getSnippetsByUserResponse struct {
	Snippets []usecases.Snippet
}

func (a *Api) endpointGetSnippetsByUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["user"], 10, 64)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	snippets, err := a.useCases.GetSnippetsByUser(usecases.UserId(userId), GetCurrentUid(r))
	if err != nil {
		WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getSnippetsByUserResponse{Snippets: snippets})
}
