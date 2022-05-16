package server

import (
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
	return echo.NewHTTPError(http.StatusInternalServerError, cv.validator.Struct(i).Error())
}

// registerAPIRoutes ...
func (h httpServer) registerAPIRoutes(routerGroup *echo.Group) {
	// Setup routes
	routerGroup.POST("/share", h.handleAPIShare)
}

// handleAPIShare ...
func (h httpServer) handleAPIShare(c echo.Context) error {
	// Set validator
	c.Echo().Validator = &CustomValidator{
		validator: validator.New(),
	}

	// Create new article share
	s := new(entity.ArticleShare)

	// Validate structure
	if err := c.Bind(s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, s)
}
