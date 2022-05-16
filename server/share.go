package server

import (
	"net/http"

	"github.com/dorneanu/gomation/server/html"
	"github.com/labstack/echo/v4"
)

// registerShareRoutes sets up routes for the share service
func (h httpServer) registerShareRoutes(routerGroup *echo.Group) {
	// Setup routes
	routerGroup.GET("/", h.handleShareIndex)
}

func (h httpServer) handleShareIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "shareIndex", html.SharePostParams{
		SendButtonMessage:   "Send article",
		CancelButtonMessage: "Cancel",
	})
}
