package v1

// Endpoint: /api/v1/snippet/{snippet}/vote
// Method: POST

import (
	"encoding/json"
	"github.com/mp-hl-2021/splinter/api"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type voteBody struct {
	Snippet usecases.SnippetId
	Vote    int
}

func (a *Api) endpointVote(w http.ResponseWriter, r *http.Request) {
	var b voteBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		api.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := a.useCases.Vote(b.Snippet, b.Vote); err != nil {
		api.WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
