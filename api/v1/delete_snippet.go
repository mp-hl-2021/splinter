package v1

// Endpoint: /api/v1/snippets/{snippet}
// Method: DELETE

import (
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/types"
	"net/http"
	"strconv"
)

func (a *Api) endpointDeleteSnippet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	snippetId, err := strconv.ParseUint(params["snippet"], 10, 64)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}
	err = a.useCases.DeleteSnippet(GetCurrentUid(r), types.SnippetId(snippetId))
	if err != nil {
		WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
