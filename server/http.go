package server

import (
	"net/http"

	"github.com/dorneanu/gomation/internal/oauth"
	"github.com/dorneanu/gomation/internal/share"
	"github.com/dorneanu/gomation/server/html"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// HTTPServerConfig holds information how to run the HTTP server
type HTTPServerConfig struct {
	ListenAddr      string
	TokenSigningKey string
	TokenExpiration int
	ShareService    share.Service
	OAuthService    oauth.Service
}

type httpServer struct {
	conf          HTTPServerConfig
	authService   oauth.Service
	shareService  share.Service
	idContextName string
}

// NewService returns a new authentication service for different providers
func NewHTTPService(s HTTPServerConfig) httpServer {
	return httpServer{
		conf:         s,
		authService:  s.OAuthService,
		shareService: s.ShareService,
		// TODO: Put this into configuration
		idContextName: "identity-provider",
	}
}

// Start starts the HTTP server
func (h httpServer) Start() {
	// New echo instance
	e := echo.New()

	// Setup middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Create index
	e.GET("/", h.handleIndex)

	// Create routing group for OAuth authentication
	authGroup := e.Group("/auth")
	h.registerAuthRoutes(authGroup)

	// Create routing group for sharing content
	shareGroup := e.Group("/share")
	h.registerShareRoutes(shareGroup)

	// Setup static content
	staticContentHandler := echo.WrapHandler(http.FileServer(http.FS(html.StaticContent)))
	e.GET("/static/*", staticContentHandler)

	// Start server
	e.Logger.Fatal(e.Start(h.conf.ListenAddr))
}

// handleIndex takes care of GET "/"
func (h httpServer) handleIndex(c echo.Context) error {
	return html.Index(c.Response().Writer, html.IndexParams{
		ProviderIndex: h.authService.ProviderIndex(),
	})
}
