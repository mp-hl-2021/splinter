package v1

// Endpoint: /api/v1/create_account
// Method: POST

import (
	"encoding/json"
	"github.com/mp-hl-2021/splinter/types"
	"net/http"
)

type createAccountBody struct {
	Username string
	Password string
}

type createAccountResponse struct {
	User types.User
}

func (a *Api) endpointCreateAccount(w http.ResponseWriter, r *http.Request) {
	var b createAccountBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	user, err := a.useCases.CreateAccount(b.Username, b.Password)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(createAccountResponse{User: user})
}
