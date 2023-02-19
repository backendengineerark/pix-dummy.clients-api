package main

import (
	"log"
	"net/http"

	"github.com/backendengineerark/clients-api/configs"
	"github.com/backendengineerark/clients-api/internal/infra/web"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	configs, err := configs.LoadConfig("./")
	if err != nil {
		panic(err)
	}

	accountHandler := web.NewAccountHandler()

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", accountHandler.CreateAccount)
	})

	log.Println("Clients API Alive!")
	http.ListenAndServe(":"+configs.AppPort, r)
}
