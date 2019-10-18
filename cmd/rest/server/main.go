package main

import (
	"log"
	"net/http"
	"time"

	"github.com/sankarvj/grpc/cmd/rest/server/internal/handlers"
)

func main() {

	api := http.Server{
		Addr:         ":8080",
		Handler:      handlers.API(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("main : API listening on %s", api.Addr)
	api.ListenAndServe()

}
