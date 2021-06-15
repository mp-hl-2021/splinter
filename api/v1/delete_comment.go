package v1

// Endpoint: /api/v1/comments/{comment}
// Method: DELETE

import (
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/types"
	"net/http"
	"strconv"
)

func (a *Api) endpointDeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	commentId, err := strconv.ParseUint(params["comment"], 10, 64)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}
	err = a.useCases.DeleteComment(GetCurrentUid(r), types.CommentId(commentId))
	if err != nil {
		WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
