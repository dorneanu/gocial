package identity

import (
	"github.com/dorneanu/gomation/internal/entity"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	// TODO: Get rid of dependency to echo.Context
	Add(entity.IdentityProvider, echo.Context) error
	GetByProvider(string, echo.Context) (entity.IdentityProvider, error)
	Save() error
	Load() error
}
