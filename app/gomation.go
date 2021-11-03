package app

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/linkedin"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

// Gomation holds anything necessary for the application to run
type Gomation struct {
	ProviderIndex ProviderIndex
}

func (g *Gomation) HandleIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("foo").Parse(indexTemplate)
	t.Execute(w, g.ProviderIndex)
}

func (g *Gomation) HandleCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	t, _ := template.New("foo").Parse(userTemplate)
	t.Execute(w, user)
}

func (g *Gomation) HandleAuth(w http.ResponseWriter, r *http.Request) {
	// try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(w, gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (g *Gomation) Init() {
	m := make(map[string]string)
	m["linkedin"] = "linkedin"
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Define here which identity providers to use
	scopes := []string{"r_emailaddress", "r_liteprofile", "w_member_social"}
	idpLinkedin := linkedin.New(
		os.Getenv("LINKEDIN_CLIENT_ID"),
		os.Getenv("LINKEDIN_CLIENT_SECRET"),
		"http://localhost:3000/auth/linkedin/callback",
		strings.Join(scopes, " "))
	goth.UseProviders(idpLinkedin)

	// Setup provider index
	g.ProviderIndex = ProviderIndex{Providers: keys, ProvidersMap: m}
}
