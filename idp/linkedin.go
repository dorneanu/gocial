package idp

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

// LinkedinProfile is used within this package as it is less useful than native types.
type LinkedinProfile struct {
	// ProfileID represents the Unique ID every Linkedin profile has.
	ProfileID string `json:"id"`
	// FirstName represents the user's first name.
	FirstName string `json:"first_name"`
	// LastName represents the user's last name.
	LastName string `json:"last-name"`
	// MaidenName represents the user's maiden name, if they have one.
	MaidenName string `json:"maiden-name"`
	// FormattedName represents the user's formatted name, based on locale.
	FormattedName string `json:"formatted-name"`
	// PhoneticFirstName represents the user's first name, spelled phonetically.
	PhoneticFirstName string `json:"phonetic-first-name"`
	// PhoneticFirstName represents the user's last name, spelled phonetically.
	PhoneticLastName string `json:"phonetic-last-name"`
	// Headline represents a short, attention grabbing description of the user.
	Headline string `json:"headline"`
	// Location represents where the user is located.
	Location Location `json:"location"`
	// Industry represents what industry the user is working in.
	Industry string `json:"industry"`
	// CurrentShare represents the user's last shared post.
	CurrentShare string `json:"current-share"`
	// NumConnections represents the user's number of connections, up to a maximum of 500.
	// The user's connections may be over this, however it will not be shown. (i.e. 500+ connections)
	NumConnections int `json:"num-connections"`
	// IsConnectionsCapped represents whether or not the user's connections are capped.
	IsConnectionsCapped bool `json:"num-connections-capped"`
	// Summary represents a long-form text description of the user's capabilities.
	Summary string `json:"summary"`
	// Specialties is a short-form text description of the user's specialties.
	Specialties string `json:"specialties"`
	// Positions is a Positions struct that describes the user's previously held positions.
	Positions Positions `json:"positions"`
	// PictureURL represents a URL pointing to the user's profile picture.
	PictureURL string `json:"picture-url"`
	// EmailAddress represents the user's e-mail address, however you must specify 'r_emailaddress'
	// to be able to retrieve this.
	EmailAddress string `json:"email-address"`
}

// Positions represents the result given by json:"positions"
type Positions struct {
	total  int
	values []Position
}

// Location specifies the users location
type Location struct {
	UserLocation string
	CountryCode  string
}

// Position represents a job held by the authorized user.
type Position struct {
	// ID represents a unique ID representing the position
	ID string
	// Title represents a user's position's title, for example Jeff Bezos's title would be 'CEO'
	Title string
	// Summary represents a short description of the user's position.
	Summary string
	// StartDate represents when the user's position started.
	StartDate string
	// EndDate represents the user's position's end date, if any.
	EndDate string
	// IsCurrent represents if the position is currently held or not.
	// If this is false, EndDate will not be returned, and will therefore equal ""
	IsCurrent bool
	// Company represents the Company where the user is employed.
	Company PositionCompany
}

// PositionCompany represents a company that is described within a user's Profile.
// This is different from Company, which fully represents a company's data.
type PositionCompany struct {
	// ID represents a unique ID representing the company
	ID string
	// Name represents the name of the company
	Name string
	// Type represents the type of the company, either 'public' or 'private'
	Type string
	// Industry represents which industry the company is in.
	Industry string
	// Ticker represents the stock market ticker symbol of the company.
	// This will be blank if the company is privately held.
	Ticker string
}

type LinkedinShare struct {
	authConf *oauth2.Config
	store    sessions.Store
	client   *http.Client
}

type LinkedinOptions struct {
	scopes       []string
	clientID     string
	clientSecret string
	redirectURL  string
}

// NewLinkedinShare returns a LinkedinShare
func NewLinkedinShare(opts LinkedinOptions) *LinkedinShare {
	_, err := url.ParseRequestURI(opts.redirectURL)
	if err != nil {
		panic(fmt.Errorf("redirectURL is invalid: %s", opts.redirectURL))
	}

	authConf := &oauth2.Config{
		ClientID:     opts.clientID,
		ClientSecret: opts.clientSecret,
		Endpoint:     linkedin.Endpoint,
		RedirectURL:  opts.redirectURL,
		Scopes:       opts.scopes,
	}

	return &LinkedinShare{
		authConf: authConf,
		store:    sessions.NewCookieStore([]byte("linkedinapi")),
	}
}

// getloginurl provides a login URL for the user to login
// TODO: Implement me
func (l *LinkedinShare) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	session, _ := l.store.Get(r, "golinkedinapi")
	session.Values["state"] = state
	defer session.Save(r, w)
}

func (l *LinkedinShare) HandleProfile(w http.ResponseWriter, r *http.Request) (*LinkedinProfile, error) {
	if l.ValidState(r) == false {
		err := fmt.Errorf("State comparison failed")
		return &LinkedinProfile{}, err
	}
	params := r.URL.Query()

	// Authenticate
	tok, err := l.authConf.Exchange(oauth2.NoContext, params.Get("code"))
	if err != nil {
		return &LinkedinProfile{}, err
	}
	client := l.authConf.Client(oauth2.NoContext, tok)
	// Retrieve data
	resp, err := l.client.Get(fullRequestURL)
	if err != nil {
		return &LinkedinProfile{}, err
	}
	// Store data to struct and return.
	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	formattedData, err := parseJSON(string(data))
	if err != nil {
		return &LinkedinProfile{}, err
	}
	return formattedData, nil
}

// generateState generates a random set of bytes to ensure state is preserved.
// This prevents such things as XSS occuring.
func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return string(b)
}

// getSessionValue grabs the value of an interface in this case being the session.Values["string"]
// This will return "" if f is nil.
func getSessionValue(f interface{}) string {
	if f != nil {
		if foo, ok := f.(string); ok {
			return foo
		}
	}
	return ""
}

// validState validates the state stored client-side with the request's state,
// it returns false if the states are not equal to each other.
func (l *LinkedinShare) ValidState(r *http.Request) bool {
	// Retrieve state
	session, _ := l.store.Get(r, "golinkedinapi")
	// Compare state to header's state
	retrievedState := session.Values["state"]
	if getSessionValue(retrievedState) != r.Header.Get("state") {
		return false
	}
	return true
}
