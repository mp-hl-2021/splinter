package highlighter

import (
	"github.com/mp-hl-2021/splinter/types"
	"log"
	"os/exec"
)

type Highlighter struct {
	requests chan types.Snippet
	storage  types.SnippetStorage
}

func MakeHighlighter(storage types.SnippetStorage) Highlighter {
	return Highlighter{
		make(chan types.Snippet),
		storage,
	}
}

func highlightSnippet(snippet *types.Snippet) (string, error) {
	out, err := exec.Command("pygmentize", "-l", string(snippet.Language), "-f", "html").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (h *Highlighter) Run() {
	for {
		snippet := <-h.requests
		hl, err := highlightSnippet(&snippet)
		if err != nil {
			log.Printf("Highlight error: %e", err)
			continue
		}
		err = h.storage.SetSnippetHighlight(snippet.Id, hl)
		if err != nil {
			log.Printf("Highlight error: %e", err)
			continue
		}
	}
}

func (h *Highlighter) Post(snippet types.Snippet) {
	h.requests <- snippet
}
