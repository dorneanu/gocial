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
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// registerAPIRoutes ...
func (h httpServer) registerAPIRoutes(routerGroup *echo.Group) {
	// Setup routes
	routerGroup.POST("/share/:provider", h.handleAPIShare)
	routerGroup.GET("/providers", h.handleAPIGetProviders)
}

// handleAPIShare ...
func (h httpServer) handleAPIShare(c echo.Context) error {
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

	// Get provider (URL parameter)
	provider := c.Param("provider")

	// Try to fetch an identity provider from the identity service
	idProvider, err := h.identityService.GetByProvider(provider, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"error":    err.Error(),
			"provider": provider,
		})
	}
	shareRepo, err := h.shareService.GetShareRepo(idProvider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"error":    err.Error(),
			"provider": provider,
		})
	}
	// Share article
	if err := h.shareService.ShareArticle(*articleShare, shareRepo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"error":    err.Error(),
			"provider": idProvider.Provider,
		})
	}
	return c.JSON(http.StatusOK, articleShare)
}

// handleAPIGetProviders ...
func (h httpServer) handleAPIGetProviders(c echo.Context) error {
	providers := make([]entity.IdentityProvider, 0)

	for _, ip := range h.availableIdentityProviders(c) {
		provider, err := h.identityService.GetByProvider(ip.Provider, c)
		if err != nil {
			continue
		}
		providers = append(providers, provider)
	}
	return c.JSONPretty(http.StatusOK, providers, "  ")
}
