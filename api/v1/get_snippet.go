package v1

// Endpoint: /api/v1/snippets/{snippet}
// Method: GET

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type getSnippetResponse struct {
	Snippet usecases.Snippet
}

func (a *Api) endpointGetSnippet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	snippetId := usecases.SnippetId(params["snippet"])

	snippet, err := a.useCases.GetSnippet(snippetId)
	if err != nil {
		WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getSnippetResponse{Snippet: snippet})
}
