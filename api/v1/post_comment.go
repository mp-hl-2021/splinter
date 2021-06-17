package v1

// Endpoint: /api/v1/snippets/{snippets}/comments
// Method: POST

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/types"
	"net/http"
	"strconv"
)

type postCommentBody struct {
	Contents string
}

type postCommentResponse struct {
	Comment types.Comment
}

func (a *Api) endpointPostComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	snippetId, err := strconv.ParseUint(params["snippet"], 10, 64)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}
	var b postCommentBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	Comment, err := a.useCases.PostComment(GetCurrentUid(r), b.Contents, types.SnippetId(snippetId))
	if err != nil {
		WriteError(w, err, http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(postCommentResponse{Comment: Comment})
}
