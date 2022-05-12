package oauth

import (
	"context"
	"net/http"

	"github.com/dorneanu/gomation/internal/entity"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

// GothRepository implements auth.Repository
type GothRepository struct {
	providerIndex entity.AuthProviderIndex
	identityRepo  entity.IdentityRepository
	jwtSigningKey string
}

func NewGothRepository(providerIndex entity.AuthProviderIndex, idRepo entity.IdentityRepository, signingKey string) *GothRepository {
	// Setup cookie store
	setupCookies()

	return &GothRepository{
		identityRepo:  idRepo,
		jwtSigningKey: signingKey,
		providerIndex: providerIndex,
	}
}

// setupCookies sets up cookies
func setupCookies() {
	// TODO: Customize this somewhere else
	key := "Secret-session-key"
	maxAge := 86400 * 30 // 30 days
	isProd := false      // set to true when serving over HTTPS

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd
	gothic.Store = store
}

// withProviderContext adds the provider parameter to go's context
func (r *GothRepository) withProviderContext(c echo.Context) *http.Request {
	// Since echo.Context is not compatible to context.Context
	// we'll add a new value to a new request add pass it to gothic
	req := c.Request().WithContext(
		context.WithValue(context.Background(), "provider", c.Param("provider")))
	return req
}

// HandleAuth does OAuth workflow for specified provider
func (r *GothRepository) HandleAuth(c echo.Context) error {
	gothic.BeginAuthHandler(c.Response(), r.withProviderContext(c))
	return nil
}

func (r *GothRepository) HandleAuthCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response().Writer, r.withProviderContext(c))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Cannot handle callback")
	}

	provider, err := gothic.GetProviderName(c.Request())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Cannot get provider")
	}

	// Create new identity provider
	id := entity.IdentityProvider{
		Provider:          provider,
		UserName:          user.Name,
		UserID:            user.UserID,
		UserDescription:   user.Description,
		UserAvatarURL:     user.AvatarURL,
		AccessToken:       user.AccessToken,
		AccessTokenSecret: user.AccessTokenSecret,
		RefreshToken:      user.RefreshToken,
		ExpiresAt:         &user.ExpiresAt,
	}

	// TODO: Put name of context variable into configuration
	// TODO: Is this secure? (passing values via context)
	c.Set("identity-provider", id)
	return nil
}
