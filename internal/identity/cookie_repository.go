package identity

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dorneanu/gomation/internal/entity"
	jwtutils "github.com/dorneanu/gomation/internal/jwt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type CookieIdentityOptions struct {
	BaseCookieName  string
	Ctx             echo.Context
	TokenSigningKey string
}

type CookieIdentityRepository struct {
	baseCookieName  string
	ctx             echo.Context
	tokenSigningKey string
}

func NewCookieIdentityRepository(opts *CookieIdentityOptions) *CookieIdentityRepository {
	return &CookieIdentityRepository{
		baseCookieName:  opts.BaseCookieName,
		ctx:             opts.Ctx,
		tokenSigningKey: opts.TokenSigningKey,
	}
}

// Add ...
func (cr *CookieIdentityRepository) Add(id entity.IdentityProvider, c echo.Context) error {
	// Generate new JWT token
	jwtToken, err := jwtutils.NewToken(id, cr.tokenSigningKey)
	if err != nil {
		return fmt.Errorf("Cannot generate new JWT token: %s", err)
	}

	identityCookie := &http.Cookie{
		Name:     fmt.Sprintf("%s-%s", cr.baseCookieName, id.Provider),
		Value:    jwtToken,
		Path:     "/",
		Expires:  time.Now().Add(36 * time.Hour),
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: 1,
	}
	c.SetCookie(identityCookie)
	return nil
}

// GetByProvider ...
func (cr *CookieIdentityRepository) GetByProvider(provider string, c echo.Context) (entity.IdentityProvider, error) {
	cookie, err := c.Cookie(fmt.Sprintf("%s-%s", cr.baseCookieName, provider))
	if err != nil {
		return entity.IdentityProvider{}, fmt.Errorf("Couldn't get cookie for provider: %s", provider)
	}

	// Parse token value
	token, err := jwt.ParseWithClaims(cookie.Value, &jwtutils.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cr.tokenSigningKey), nil
	})

	// Check if valid
	if claims, ok := token.Claims.(*jwtutils.JwtCustomClaims); ok && token.Valid {
		return entity.IdentityProvider{
			Provider:          claims.Provider,
			UserName:          claims.UserName,
			UserID:            claims.UserID,
			UserDescription:   claims.UserDescription,
			UserAvatarURL:     claims.UserAvatarURL,
			AccessToken:       claims.AccessToken,
			AccessTokenSecret: claims.AccessTokenSecret,
			RefreshToken:      claims.RefreshToken,
			ExpiresAt:         &time.Time{},
		}, nil
	} else {
		return entity.IdentityProvider{}, fmt.Errorf("Couldn't validate JWT token: %s", err)
	}

}

// Load
func (cr *CookieIdentityRepository) Load() error {
	return nil
}

// Save ...
func (cr *CookieIdentityRepository) Save() error {
	return nil
}
