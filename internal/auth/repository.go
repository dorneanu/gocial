package auth

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/linkedin"
)

type Repository interface {
	Init()
	HandleIndex(http.ResponseWriter, *http.Request)
	HandleAuth(http.ResponseWriter, *http.Request)
	HandleCallback(http.ResponseWriter, *http.Request)
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

type GothRepository struct {
	mux           *mux.Router
	providerIndex ProviderIndex
	oauthConfigs  []OAuthConfig
}

func NewGothRepository(m *mux.Router, confs []OAuthConfig) *GothRepository {
	return &GothRepository{
		mux:          m,
		oauthConfigs: confs,
	}
}

func (r *GothRepository) Init() {
	r.setupProviders()
	r.setupCookies()
	r.setupRoutes()
}

func (r *GothRepository) setupCookies() {
	// TODO: Replace this
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

func (r *GothRepository) setupRoutes() {
	r.mux.HandleFunc("/", r.HandleIndex)
	r.mux.HandleFunc("/auth/{provider}", r.HandleAuth)
	r.mux.HandleFunc("/auth/{provider}/callback", r.HandleCallback)
}

func (r *GothRepository) setupProviders() {
	m := make(map[string]string)

	for _, oauthConf := range r.oauthConfigs {
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
			fmt.Printf("Needs to be implemented")
		}
	}

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	r.providerIndex = ProviderIndex{Providers: keys, ProvidersMap: m}
	sort.Strings(keys)
}

func (r *GothRepository) HandleIndex(w http.ResponseWriter, req *http.Request) {
	t, _ := template.New("gothrepository").Parse(indexTemplate)
	t.Execute(w, r.providerIndex)
}

func (r *GothRepository) HandleAuth(w http.ResponseWriter, req *http.Request) {
	gothic.BeginAuthHandler(w, req)
}

func (r *GothRepository) HandleCallback(w http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(w, req)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	t, _ := template.New("gothrepository").Parse(userTemplate)
	t.Execute(w, user)
}
