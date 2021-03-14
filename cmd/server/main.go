package main

import (
	"github.com/mp-hl-2021/splinter/internal/app/splinter/api"
	"github.com/mp-hl-2021/splinter/internal/app/splinter/backend"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", api.NewApi(backend.NewSimpleUserInterface()).Router())
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Panic(err)
	}
}
