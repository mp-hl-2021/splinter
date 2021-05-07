package v1

// Endpoint: /api/v1/authenticate
// Method: POST

import (
	"encoding/json"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type authenticateBody struct {
	Username string
	Password string
}

type authenticateResponse struct {
	Token usecases.Token
}

func (a *Api) endpointAuthenticate(w http.ResponseWriter, r *http.Request) {
	var b authenticateBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	token, err := a.useCases.Authenticate(b.Username, b.Password)
	if err != nil {
		WriteError(w, err, http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(authenticateResponse{Token: token})
}
