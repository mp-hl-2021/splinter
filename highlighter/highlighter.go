package highlighter

import (
	"errors"
	"github.com/mp-hl-2021/splinter/types"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var ChannelIsFullErr = errors.New("channel is full")

type Highlighter struct {
	requests chan types.Snippet
	storage  types.SnippetStorage
}

func New(storage types.SnippetStorage, size int) Highlighter {
	return Highlighter{
		make(chan types.Snippet, size),
		storage,
	}
}

func highlightSnippet(snippet *types.Snippet) (string, error) {
	input, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer os.Remove(input.Name())
	defer input.Close()

	if _, err = input.WriteString(snippet.Contents); err != nil {
		return "", err
	}

	cmd := exec.Command("pygmentize", "-l", string(snippet.Language), "-f", "html", input.Name())
	out, err := cmd.Output()
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

func (h *Highlighter) Post(snippet types.Snippet) error {
	select {
	case h.requests <- snippet:
		return nil
	default:
		return ChannelIsFullErr
	}
}
