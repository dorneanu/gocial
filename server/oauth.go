package server

import (
	"fmt"
	"net/http"

	"github.com/dorneanu/gocial/internal/entity"
	"github.com/dorneanu/gocial/server/html"
	"github.com/labstack/echo/v4"
)

func (h httpServer) registerAuthRoutes(routerGroup *echo.Group) {
	// setup
	// jwtConfig := middleware.JWTConfig{
	// 	Claims:      &jwtutils.JwtCustomClaims{},
	// 	SigningKey:  []byte(h.conf.TokenSigningKey),
	// 	TokenLookup: "cookie:gocialAuth",
	// }

	// Setup routes
	routerGroup.GET("/", h.handleOAuthIndex)
	routerGroup.GET("/info", h.handleOAuthInfo)
	routerGroup.GET("/logout", h.handleOAuthLogout)
	routerGroup.GET("/:provider", h.handleOAuth)
	routerGroup.GET("/callback/:provider", h.handleOAuthCallback)
	// routerGroup.GET("/info",
	// 	h.handleOAuthInfo,
	// 	middleware.JWTWithConfig(jwtConfig),
	// )
	routerGroup.GET("/info", h.handleOAuthInfo)
}

// handleOAuthIndex handles index page for OAuth workflow
func (h httpServer) handleOAuthIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "authIndex", html.AuthIndexParams{
		ProviderIndex: h.authService.ProviderIndex(),
	})
}

// handleOAuthInfo shows information about current authentications
func (h httpServer) handleOAuthInfo(c echo.Context) error {
	return c.Render(http.StatusOK, "authInfo", h.availableIdentityProviders(c))
}

// handleOAuth handles OAuth workflow
func (h httpServer) handleOAuth(c echo.Context) error {
	return h.authService.Repo().HandleAuth(c)
}

// handleOAuthCallback ...
func (h httpServer) handleOAuthCallback(c echo.Context) error {
	err := h.authService.Repo().HandleAuthCallback(c)
	if err != nil {
		return err
	}

	// Fetch new identity provider
	identityProvider := c.Get(h.idContextName).(entity.IdentityProvider)

	// Persis new identity provider
	h.conf.IdentityService.Add(identityProvider, c)

	// TODO: Put /auth/info into configuration
	return c.Redirect(http.StatusTemporaryRedirect, "/auth/info")
}

// availableIdentityProviders returns a list of all available identity providers
func (h httpServer) availableIdentityProviders(c echo.Context) []entity.IdentityProvider {
	var identityProviders []entity.IdentityProvider

	for _, p := range h.providerIndex.Providers {
		// Try to fetch an identity provider from the identity service
		idProvider, err := h.identityService.GetByProvider(p, c)
		if err != nil {
			fmt.Printf("Provider %s not found\n", p)
			continue
		}
		identityProviders = append(identityProviders, idProvider)
	}
	return identityProviders
}

// handleOAuthLogout ...
func (h httpServer) handleOAuthLogout(c echo.Context) error {
	// Delete all OAuth related cookies
	for _, p := range h.providerIndex.Providers {
		h.identityService.Delete(p, c)
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
