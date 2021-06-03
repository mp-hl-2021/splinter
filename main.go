package main

import (
	"flag"
	"github.com/mp-hl-2021/splinter/api"
	"github.com/mp-hl-2021/splinter/auth"
	"github.com/mp-hl-2021/splinter/storage"
	"github.com/mp-hl-2021/splinter/usecases"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
	publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
	connStr := flag.String("connStr", "user=postgres password=postgres host=db dbname=postgres sslmode=disable", "postgres connection string")
	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyPath)
	publicKeyBytes, err := ioutil.ReadFile(*publicKeyPath)

	a, err := auth.NewJwtHandler(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}

	postgres, err := storage.NewPostgres(*connStr)
	if err != nil {
		panic(err)
	}

	userInterface := &usecases.DelegatedUserInterface{
		UserStorage:    postgres,
		SnippetStorage: postgres,
		Auth:           a,
	}

	service := api.NewApi(userInterface)
	addr := ":5000"

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
