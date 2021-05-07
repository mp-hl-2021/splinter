package main

import (
	"github.com/mp-hl-2021/splinter/api"
	"github.com/mp-hl-2021/splinter/usecases"
	"log"
	"net/http"
	"time"
)

func main() {
	service := api.NewApi(usecases.DummyUserInterface{})
	addr := "127.0.0.1:5000"

	server := http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}

	log.Printf("Serving at %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
