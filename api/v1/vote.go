package v1

// Endpoint: /api/v1/snippet/{snippet}/vote
// Method: POST

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/types"
	"net/http"
	"strconv"
)

type voteBody struct {
	Vote    int
}

func (a *Api) endpointVote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var b voteBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	snippetId, err := strconv.ParseUint(params["snippet"], 10, 64)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := a.useCases.Vote(GetCurrentUid(r), types.SnippetId(snippetId), b.Vote); err != nil {
		WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
