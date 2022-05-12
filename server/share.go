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
	routerGroup.POST("/article", h.handleShareArticle)
	routerGroup.POST("/comment", h.handleShareComment)
}

func (h httpServer) handleShareIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "post", html.PostParams{
		SendButtonMessage:   "Send article",
		CancelButtonMessage: "Cancel",
	})
}

// handleShareArticle posts a new article to different providers
func (h httpServer) handleShareArticle(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented yet")
}

// handleShareComment posts a new comment to different providers
func (h httpServer) handleShareComment(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented yet")
}
