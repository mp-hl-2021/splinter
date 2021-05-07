package v1

// Endpoint: /api/v1/users/current
// Method: GET

import (
	"encoding/json"
	"github.com/mp-hl-2021/splinter/api"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type getCurrentUserResponse struct {
	User usecases.User
}

func (a *Api) endpointGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := a.useCases.GetCurrentUser()
	if err != nil {
		api.WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getCurrentUserResponse{User: user})
}
