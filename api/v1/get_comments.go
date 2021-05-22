package v1

// Endpoint: /snippets/{snippet}/comments
// Method: GET

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
	"strconv"
)

type getCommentsResponse struct {
	Comments []usecases.Comment
}

func (a *Api) endpointGetComments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	snippetId, err := strconv.ParseUint(params["snippet"], 10, 64)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}
	comments, err := a.useCases.GetComments(usecases.SnippetId(snippetId))
	if err != nil {
		WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getCommentsResponse{Comments: comments})
}
