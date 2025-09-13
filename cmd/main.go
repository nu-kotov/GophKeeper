package main

import (
	"log"

	"github.com/go-chi/chi/v5"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	r := chi.NewRouter()
}
