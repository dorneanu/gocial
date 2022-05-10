package oauth

import "github.com/labstack/echo/v4"

type Repository interface {
	HandleAuth(echo.Context) error
	HandleAuthCallback(echo.Context) error
}
