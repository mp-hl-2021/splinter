package v1

// Endpoint: /api/v1/snippets/{snippets}/comments
// Method: POST

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type postCommentBody struct {
	Contents string
}

type postCommentResponse struct {
	Comment usecases.Comment
}

func (a *Api) endpointPostComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	snippetId := usecases.SnippetId(params["snippet"])

	var b postCommentBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	Comment, err := a.useCases.PostComment(b.Contents, snippetId)
	if err != nil {
		WriteError(w, err, http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(postCommentResponse{Comment: Comment})
}
