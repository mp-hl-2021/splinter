package main

import (
	"flag"
	"github.com/mp-hl-2021/splinter/api"
	"github.com/mp-hl-2021/splinter/auth"
	"github.com/mp-hl-2021/splinter/highlighter"
	"github.com/mp-hl-2021/splinter/storage"
	"github.com/mp-hl-2021/splinter/usecases"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
	publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
	connStr := flag.String("connStr", "user=postgres password=postgres host=db dbname=postgres sslmode=disable", "postgres connection string")
	highlightWorkers := flag.String("highlightWorkers", "8", "number of highlighter workers")
	highlightQueueSize := flag.String("highlightQueueSize", "256", "highlighter queue size")
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

	hq, err := strconv.Atoi(*highlightQueueSize)
	if err != nil {
		panic(err)
	}

	h := highlighter.New(postgres, hq)

	w, err := strconv.Atoi(*highlightWorkers)
	if err != nil {
		panic(err)
	}

	for i := 0; i < w; i += 1 {
		go h.Run()
	}

	userInterface := &usecases.DelegatedUserInterface{
		UserStorage:    postgres,
		SnippetStorage: postgres,
		Auth:           a,
		Highlighter:    h,
	}

	service := api.NewApi(userInterface, a)
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
