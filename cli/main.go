package main

import (
	"github.com/dorneanu/gomation/idp"
	"github.com/gorilla/mux"
)

func main() {
	opts := idp.LinkedinOptions{
		scopes:       []string{"bla"},
		clientID:     "boru",
		clientSecret: "alles",
		redirectURL:  "http://localhost:1313",
	}
	linkedinShare := idp.NewLinkedinShare(opts)
	router := mux.NewRouter()
	router.HandleFunc("/login", linkedinShare.HandleLogin)
}
