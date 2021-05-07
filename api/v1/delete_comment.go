package v1

// Endpoint: /api/v1/comments/{comment}
// Method: DELETE

import (
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

func (a *Api) endpointDeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	commentId := usecases.CommentId(params["comment"])

	err := a.useCases.DeleteComment(commentId)
	if err != nil {
		WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
