package server

import (
	"html/template"
	"io"
	"net/http"

	"github.com/dorneanu/gomation/internal/entity"
	"github.com/dorneanu/gomation/internal/identity"
	"github.com/dorneanu/gomation/internal/oauth"
	"github.com/dorneanu/gomation/internal/share"
	"github.com/dorneanu/gomation/server/html"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// HTTPServerConfig holds information how to run the HTTP server
type HTTPServerConfig struct {
	ListenAddr      string
	TokenSigningKey string
	TokenExpiration int
	ShareService    share.Service
	OAuthService    oauth.Service
	IdentityService identity.Repository
	ProviderIndex   *entity.AuthProviderIndex
}

type httpServer struct {
	ctx             echo.Context
	conf            HTTPServerConfig
	authService     oauth.Service
	shareService    share.Service
	identityService identity.Repository
	providerIndex   *entity.AuthProviderIndex
	idContextName   string
}

// NewService returns a new authentication service for different providers
func NewHTTPService(s HTTPServerConfig) httpServer {
	return httpServer{
		conf:            s,
		authService:     s.OAuthService,
		shareService:    s.ShareService,
		identityService: s.IdentityService,
		providerIndex:   s.ProviderIndex,
		// TODO: Put this into configuration
		idContextName: "identity-provider",
	}
}

// Start starts the HTTP server
func (h httpServer) Start(e *echo.Echo) {

	// Setup middleware
	e.Use(middleware.Logger())
	e.Logger.SetLevel(99)
	// e.Debug = true
	e.Use(middleware.Recover())

	// Setup HTML templating
	e.Renderer = html.RegisterTemplates()

	// Create general routes
	e.GET("/", h.handleIndex)
	e.GET("/about/", h.handleAbout)

	// Create routing group for OAuth authentication
	authGroup := e.Group("/auth")
	h.registerAuthRoutes(authGroup)

	// Create routing group for sharing content
	shareGroup := e.Group("/share")
	h.registerShareRoutes(shareGroup)

	// Create routing group for the REST API
	apiGroup := e.Group("/api")
	h.registerAPIRoutes(apiGroup)

	// Setup static content
	staticContentHandler := echo.WrapHandler(http.FileServer(http.FS(html.StaticContent)))
	e.GET("/static/*", staticContentHandler)
}

// handleIndex takes care of GET "/"
func (h httpServer) handleIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (h httpServer) handleAbout(c echo.Context) error {
	return c.Render(http.StatusOK, "about", nil)
}
