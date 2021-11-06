package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type OAuthConfig struct {
	ProviderName string
	ClientID     string
	ClientSecret string
	CallbackURL  string
	Scopes       []string
}

type Service interface {
	Start()
}

type service struct {
	repo   Repository
	router *mux.Router
}

func NewService(repo Repository, router *mux.Router) Service {
	return service{
		repo:   repo,
		router: router,
	}
}

func (s service) Start() {
	s.repo.Init()

	// New web server
	// TODO: Put web server configs into some user config
	srv := &http.Server{
		Handler:      s.router,
		Addr:         "127.0.0.1:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
