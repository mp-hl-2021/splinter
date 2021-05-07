package v1

// Endpoint: /api/v1/snippets/{snippet}
// Method: DELETE

import (
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/api"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

func (a *Api) endpointDeleteSnippet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	snippetId := usecases.SnippetId(params["snippet"])

	err := a.useCases.DeleteSnippet(snippetId)
	if err != nil {
		api.WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
