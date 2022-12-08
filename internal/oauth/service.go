package oauth

import (
	"strings"

	"github.com/dorneanu/gocial/internal/entity"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/twitter"
)

type OAuthConfig struct {
	ProviderName     string
	ClientID         string
	ClientSecret     string
	CallbackURL      string
	Scopes           []string
	IdentityProvider entity.IdentityProvider
}

type Service interface {
	Repo() Repository
	ProviderIndex() entity.AuthProviderIndex
}

type ServiceConfig struct {
	Repo          Repository
	ProviderIndex entity.AuthProviderIndex
}

// oauthService implements oauth.Service (interface)
type oauthService struct {
	repo          Repository
	providerIndex entity.AuthProviderIndex
}

// NewService returns a new authentication service for different providers
func NewService(conf ServiceConfig) Service {
	return oauthService{
		repo:          conf.Repo,
		providerIndex: conf.ProviderIndex,
	}
}

// Repo() returns repository which implements oauth.Repository
func (s oauthService) Repo() Repository {
	return s.repo
}

// ProviderIndex returns an index to all available authentication providers
func (s oauthService) ProviderIndex() entity.AuthProviderIndex {
	return s.providerIndex
}

// SetupAuthProviders configures lits of supported authentication providers
func SetupAuthProviders(confs []OAuthConfig) entity.AuthProviderIndex {
	oauthConfs := make(map[string]OAuthConfig)
	for _, c := range confs {
		oauthConfs[c.ProviderName] = c
	}

	m := make(map[string]string)
	for _, oauthConf := range oauthConfs {
		m[oauthConf.ProviderName] = oauthConf.ProviderName

		// TODO: Change this to switch case?
		if oauthConf.ProviderName == "linkedin" {
			idpLinkedin := linkedin.New(
				oauthConf.ClientID,
				oauthConf.ClientSecret,
				oauthConf.CallbackURL,
				strings.Join(oauthConf.Scopes, " "))

			goth.UseProviders(idpLinkedin)
		} else if oauthConf.ProviderName == "twitter" {
			idpTwitter := twitter.New(
				oauthConf.ClientID,
				oauthConf.ClientSecret,
				oauthConf.CallbackURL,
			)
			goth.UseProviders(idpTwitter)
		}
	}

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return entity.AuthProviderIndex{Providers: keys, ProvidersMap: m}
}
