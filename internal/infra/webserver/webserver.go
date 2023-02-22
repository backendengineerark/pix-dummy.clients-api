package webserver

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/backendengineerark/clients-api/internal/infra/database"
	"github.com/backendengineerark/clients-api/internal/infra/webserver/custom_middleware"
	"github.com/backendengineerark/clients-api/internal/infra/webserver/handlers"
	"github.com/backendengineerark/clients-api/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	WebServerPort int
	Db            *sql.DB
}

func NewWebServer(webServerPort int, db *sql.DB) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		WebServerPort: webServerPort,
		Db:            db,
	}
}

func (ws *WebServer) Start() {
	ws.Router.Use(middleware.Recoverer)
	ws.Router.Use(custom_middleware.LoggerWithCorrelationId)

	clientRepository := database.NewClientRepository(ws.Db)
	accountRepository := database.NewAccountRepository(ws.Db)
	createAccountUseCase := usecase.NewCreateAccountUseCase(*ws.Db, clientRepository, accountRepository)
	accountHandler := handlers.NewAccountHandler(createAccountUseCase)

	ws.Router.Route("/accounts", func(r chi.Router) {
		r.Post("/", accountHandler.CreateAccount)
	})

	log.Println("WebServer up!")
	http.ListenAndServe(":"+strconv.Itoa(ws.WebServerPort), ws.Router)
}
