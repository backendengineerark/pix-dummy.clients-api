package webserver

import (
	"log"
	"net/http"
	"strconv"

	"github.com/backendengineerark/clients-api/internal/infra/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	WebServerPort int
}

func NewWebServer(webServerPort int) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		WebServerPort: webServerPort,
	}
}

func (ws *WebServer) Start() {
	ws.Router.Use(middleware.Logger)
	ws.Router.Use(middleware.Recoverer)

	accountHandler := handlers.NewAccountHandler()

	ws.Router.Route("/accounts", func(r chi.Router) {
		r.Post("/", accountHandler.CreateAccount)
	})

	log.Println("WebServer up!")
	http.ListenAndServe(":"+strconv.Itoa(ws.WebServerPort), ws.Router)
}
