package server

import (
	"net/http"
	"time"

	"github.com/dorneanu/gomation/internal/entity"
	jwtutils "github.com/dorneanu/gomation/internal/jwt"
	"github.com/dorneanu/gomation/server/html"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (h httpServer) registerAuthRoutes(routerGroup *echo.Group) {
	// setup
	jwtConfig := middleware.JWTConfig{
		Claims:      &jwtutils.JwtCustomClaims{},
		SigningKey:  []byte(h.conf.TokenSigningKey),
		TokenLookup: "cookie:gocialAuth",
	}

	// Setup routes
	routerGroup.GET("/", h.handleOAuthIndex)
	routerGroup.GET("/info", h.handleOAuthInfo)
	routerGroup.GET("/:provider", h.handleOAuth)
	routerGroup.GET("/callback/:provider", h.handleOAuthCallback)
	routerGroup.GET("/info",
		h.handleOAuthInfo,
		middleware.JWTWithConfig(jwtConfig),
	)
}

// handleOAuthIndex handles index page for OAuth workflow
func (h httpServer) handleOAuthIndex(c echo.Context) error {
	return html.Index(c.Response().Writer, html.IndexParams{
		ProviderIndex: h.authService.ProviderIndex(),
	})
}

// handleOAuthInfo shows information about current authentications
func (h httpServer) handleOAuthInfo(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtutils.JwtCustomClaims)

	return html.Profile(c.Response().Writer, *claims)
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

	// Fetch new identity probvider
	identProvider := c.Get(h.idContextName).(entity.IdentityProvider)

	// Create new JWT token
	jwtToken, err := jwtutils.NewToken(identProvider, h.conf.TokenSigningKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Cannot generate new JWT token",
			"error":   err,
		})
	}

	// // Send cookie back
	authCookie := &http.Cookie{
		Name:     "gocialAuth",
		Value:    jwtToken,
		Path:     "/",
		Expires:  time.Now().Add(36 * time.Hour),
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: 1,
	}
	c.SetCookie(authCookie)

	// TODO: Put /auth/info into configuration
	return c.Redirect(http.StatusTemporaryRedirect, "/auth/info")
}
