package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	myRouter := chi.NewRouter()
	myRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})

	http.ListenAndServe(":8080", myRouter)
}
