package v1

// Endpoint: /api/v1/snippets/{snippet}
// Method: GET

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
	"strconv"
)

type getSnippetResponse struct {
	Snippet usecases.Snippet
}

func (a *Api) endpointGetSnippet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	snippetId, err := strconv.ParseUint(params["snippet"], 10, 64)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	snippet, err := a.useCases.GetSnippet(GetCurrentUid(r), usecases.SnippetId(snippetId))
	if err != nil {
		WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getSnippetResponse{Snippet: snippet})
}
