package auth

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/dorneanu/gomation/internal/entity"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/twitter"
)

type Repository interface {
	Init()
	HandleIndex(http.ResponseWriter, *http.Request)
	HandleAuth(http.ResponseWriter, *http.Request)
	HandleCallback(http.ResponseWriter, *http.Request)
	SaveAuthDetails(http.ResponseWriter, *http.Request)
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

type GothRepository struct {
	mux           *mux.Router
	providerIndex ProviderIndex
	oauthConfigs  map[string]OAuthConfig
	identityRepo  entity.IdentityRepository
}

func NewGothRepository(m *mux.Router, confs []OAuthConfig, idRepo entity.IdentityRepository) *GothRepository {
	oauthConfs := make(map[string]OAuthConfig)
	for _, c := range confs {
		oauthConfs[c.ProviderName] = c
	}

	return &GothRepository{
		mux:          m,
		oauthConfigs: oauthConfs,
		identityRepo: idRepo,
	}
}

func (r *GothRepository) Init() {
	r.setupProviders()
	r.setupCookies()
	r.setupRoutes()
}

func (r *GothRepository) setupCookies() {
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

func (r *GothRepository) setupRoutes() {
	// Setup HTTP routes using gorilla/mux
	r.mux.HandleFunc("/", r.HandleIndex)
	r.mux.HandleFunc("/auth/{provider}", r.HandleAuth)
	r.mux.HandleFunc("/auth/{provider}/callback", r.HandleCallback)
	r.mux.HandleFunc("/save", r.SaveAuthDetails)
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

	provider, err := gothic.GetProviderName(req)
	if err != nil {
		fmt.Fprintf(w, "Couldn't get provider: %s", err)
		return
	}

	// Create new identity
	id := entity.Identity{
		Provider:          provider,
		Name:              user.Name,
		ID:                user.UserID,
		AccessToken:       user.AccessToken,
		AccessTokenSecret: user.AccessTokenSecret,
	}

	fmt.Printf("repo: %v\n", id)

	// Persist identity
	err = r.identityRepo.Add(id)
	if err != nil {
		fmt.Printf("Couldn't add new identity: %s\n", err)
	}

	// Redirect to save
	http.Redirect(w, req, "/save", http.StatusSeeOther)
}

func (r *GothRepository) SaveAuthDetails(w http.ResponseWriter, req *http.Request) {
	err := r.identityRepo.Save()
	if err != nil {
		fmt.Fprintf(w, "There was an error saving auth details: %s\n", err)
	} else {
		fmt.Fprint(w, "Successfully stored auth details")
	}
}
