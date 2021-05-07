package v1

// Endpoint: /snippets/{snippet}/comments
// Method: GET

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/api"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type getCommentsResponse struct {
	Comments []usecases.Comment
}

func (a *Api) endpointGetComments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	snippetId := usecases.SnippetId(params["snippet"])

	comments, err := a.useCases.GetComments(snippetId)
	if err != nil {
		api.WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getCommentsResponse{Comments: comments})
}
