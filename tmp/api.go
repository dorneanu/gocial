package golinkedinapi

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

const (
	fullRequestURL  = "https://api.linkedin.com/v1/people/~:(id,first-name,email-address,last-name,headline,picture-url,industry,summary,specialties,positions:(id,title,summary,start-date,end-date,is-current,company:(id,name,type,size,industry,ticker)),educations:(id,school-name,field-of-study,start-date,end-date,degree,activities,notes),associations,interests,num-recommenders,date-of-birth,publications:(id,title,publisher:(name),authors:(id,name),date,url,summary),patents:(id,title,summary,number,status:(id,name),office:(name),inventors:(id,name),date,url),languages:(id,language:(name),proficiency:(level,name)),skills:(id,skill:(name)),certifications:(id,name,authority:(name),number,start-date,end-date),courses:(id,name,number),recommendations-received:(id,recommendation-type,recommendation-text,recommender),honors-awards,three-current-positions,three-past-positions,volunteer)?format=json"
	basicRequestURL = "https://api.linkedin.com/v1/people/~:(id,first-name,last-name,headline,picture-url,industry,summary,specialties,positions:(id,title,summary,start-date,end-date,is-current,company:(id,name,type,size,industry,ticker)),educations:(id,school-name,field-of-study,start-date,end-date,degree,activities,notes),associations,interests,num-recommenders,date-of-birth,publications:(id,title,publisher:(name),authors:(id,name),date,url,summary),patents:(id,title,summary,number,status:(id,name),office:(name),inventors:(id,name),date,url),languages:(id,language:(name),proficiency:(level,name)),skills:(id,skill:(name)),certifications:(id,name,authority:(name),number,start-date,end-date),courses:(id,name,number),recommendations-received:(id,recommendation-type,recommendation-text,recommender),honors-awards,three-current-positions,three-past-positions,volunteer)?format=json"
)

var (
	validPermissions = map[string]bool{
		"r_basicprofile":   true,
		"r_emailaddress":   true,
		"rw_company_admin": true,
		"w_share":          true,
	}
	authConf     *oauth2.Config
	store        = sessions.NewCookieStore([]byte("golinkedinapi"))
	requestedURL string
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

// ParseJSON converts a JSON string to a pointer to a LinkedinProfile.
func parseJSON(s string) (*LinkedinProfile, error) {
	linkedinProfile := &LinkedinProfile{}
	bytes := bytes.NewBuffer([]byte(s))
	err := json.NewDecoder(bytes).Decode(linkedinProfile)
	if err != nil {
		return nil, err
	}
	return linkedinProfile, nil
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
func validState(r *http.Request) bool {
	// Retrieve state
	session, _ := store.Get(r, "golinkedinapi")
	// Compare state to header's state
	retrievedState := session.Values["state"]
	if getSessionValue(retrievedState) != r.Header.Get("state") {
		return false
	}
	return true
}

// GetLoginURL provides a state-specific login URL for the user to login to.
// CAUTION: This must be called before GetProfileData() as this enforces state.
func GetLoginURL(w http.ResponseWriter, r *http.Request) string {
	state := generateState()
	// The only time this will error out is if the session cannot be decoded, however
	// we don't care about that as we can simply change state.
	session, _ := store.Get(r, "golinkedinapi")
	session.Values["state"] = state
	defer session.Save(r, w)
	return authConf.AuthCodeURL(state)
}

// GetProfileData gather's the user's Linkedin profile data and returns it as a pointer to a LinkedinProfile struct.
// CAUTION: GetLoginURL must be called before this, as GetProfileData() has a state check.
func GetProfileData(w http.ResponseWriter, r *http.Request) (*LinkedinProfile, error) {
	if validState(r) == false {
		err := fmt.Errorf("State comparison failed")
		return &LinkedinProfile{}, err
	}
	params := r.URL.Query()
	// Authenticate
	tok, err := authConf.Exchange(oauth2.NoContext, params.Get("code"))
	if err != nil {
		return &LinkedinProfile{}, err
	}
	client := authConf.Client(oauth2.NoContext, tok)
	// Retrieve data
	resp, err := client.Get(fullRequestURL)
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

// InitConfig initializes the config needed by the client.
// permissions is a string of all scopes desired by the user.
func InitConfig(permissions []string, clientID string, clientSecret string, redirectURL string) {
	var isEmail, isBasic bool
	if permissions == nil {
		panic(fmt.Errorf("You must specify some scope to request"))
	}
	for _, elem := range permissions {
		if elem == "r_emailaddress" {
			isEmail = true
		} else if elem == "r_basicprofile" {
			isBasic = true
		}
		if validPermissions[elem] != true {
			panic(fmt.Errorf("All elements of permissions must be valid Linkedin permissions as specified in the API docs"))
		}
	}
	_, err := url.ParseRequestURI(redirectURL)
	if err != nil {
		panic(fmt.Errorf("redirectURL specified must be a valid FQDN. Please ensure you added https:// to the front"))
	}
	authConf = &oauth2.Config{ClientID: clientID,
		ClientSecret: clientSecret,
		Endpoint:     linkedin.Endpoint,
		RedirectURL:  redirectURL,
		Scopes:       permissions,
	}
	if isEmail && isBasic {
		requestedURL = fullRequestURL
	} else if isBasic {
		requestedURL = basicRequestURL
	}
}

// TODO: Config Struct
