package main

import (
	"os"

	"github.com/dorneanu/gomation/internal/auth"
	"github.com/gorilla/mux"
)

func main() {
	// New linkedin repository
	oauthConfigs := []auth.OAuthConfig{
		auth.OAuthConfig{
			ProviderName: "linkedin",
			Scopes:       []string{"r_emailaddress", "r_liteprofile", "w_member_social"},
			ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
			ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
			CallbackURL:  "http://localhost:3000/auth/linkedin/callback",
		},
	}
	// New goth auth repository
	r := mux.NewRouter()
	gothRepository := auth.NewGothRepository(r, oauthConfigs)

	// New auth service
	authService := auth.NewService(gothRepository, r)
	authService.Start()
}
