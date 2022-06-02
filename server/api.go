package server

import (
	"fmt"
	"net/http"

	"github.com/dorneanu/gomation/internal/entity"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator validates input from POST request
//
// Found at https://github.com/labstack/echo/issues/1803
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// registerAPIRoutes ...
func (h httpServer) registerAPIRoutes(routerGroup *echo.Group) {
	// Setup routes
	routerGroup.POST("/share", h.handleAPIShare)
}

// handleAPIShare ...
func (h httpServer) handleAPIShare(c echo.Context) error {
	fmt.Print("API Share")

	// Custom validator
	c.Echo().Validator = &CustomValidator{validator: validator.New()}

	// Create new article share
	articleShare := new(entity.ArticleShare)

	// Validate structure
	if err := c.Bind(articleShare); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(articleShare); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get available identity providers
	for _, ip := range h.availableIdentityProviders(c) {
		shareRepo, err := h.shareService.GetShareRepo(ip)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		// Share article
		if err := h.shareService.ShareArticle(*articleShare, shareRepo); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
				"error":    err.Error(),
				"provider": ip.Provider,
			})
		}
	}

	return c.JSON(http.StatusOK, articleShare)
}
