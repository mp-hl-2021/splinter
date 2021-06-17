package v1

// Endpoint: /api/v1/snippets
// Method: POST

import (
	"encoding/json"
	"github.com/mp-hl-2021/splinter/types"
	"net/http"
)

type postSnippetBody struct {
	Contents string
	Language types.ProgrammingLanguage
}

type postSnippetResponse struct {
	Snippet types.Snippet
}

func (a *Api) endpointPostSnippet(w http.ResponseWriter, r *http.Request) {
	var b postSnippetBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	snippet, err := a.useCases.PostSnippet(GetCurrentUid(r), b.Contents, b.Language)
	if err != nil {
		WriteError(w, err, http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(postSnippetResponse{Snippet: snippet})
}
