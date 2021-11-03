package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dorneanu/gomation/app"
	"github.com/gorilla/mux"
)

func main() {
	// Create new application
	app := &app.Gomation{}
	app.Init()

	// Create http mux
	mux := mux.NewRouter()
	mux.HandleFunc("/", app.HandleIndex)
	mux.HandleFunc("/auth/{provider}", app.HandleAuth)
	mux.HandleFunc("/auth/{provider}/callback", app.HandleCallback)

	// new HTTP server
	srv := &http.Server{
		Handler:      mux,
		Addr:         "127.0.0.1:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
