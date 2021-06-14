package v1

// Endpoint: /api/v1/snippets/language/{language}
// Method: GET

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/splinter/usecases"
	"net/http"
)

type getSnippetsByLanguageResponse struct {
	Snippets []usecases.Snippet
}

func (a *Api) endpointGetSnippetsByLanguage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	language := usecases.ProgrammingLanguage(params["language"])

	snippets, err := a.useCases.GetSnippetsByLanguage(language, GetCurrentUid(r))
	if err != nil {
		WriteError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getSnippetsByLanguageResponse{Snippets: snippets})
}
