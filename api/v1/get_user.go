package v1

// Endpoint: /api/v1/users/{user}
// Method: GET

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/types"
	"net/http"
	"strconv"
)

type getUserResponse struct {
	User types.User
}

func (a *Api) endpointGetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["user"], 10, 64)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}
	user, err := a.useCases.GetUser(types.UserId(userId))
	if err != nil {
		WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getUserResponse{User: user})
}
