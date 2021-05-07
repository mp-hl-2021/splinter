package v1

// Endpoint: /api/v1/users/{user}
// Method: GET

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/api"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type getUserResponse struct {
	User usecases.User
}

func (a *Api) endpointGetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := usecases.UserId(params["user"])

	user, err := a.useCases.GetUser(userId)
	if err != nil {
		api.WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getUserResponse{User: user})
}
