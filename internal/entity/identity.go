package entity

// Identity holds information about an identity
type Identity struct {
	Provider          string `yaml:"provider"`
	Name              string `yaml:"name"`
	ID                string `yaml:"id"`
	AccessToken       string `yaml:"accessToken"`
	AccessTokenSecret string `yaml:"accessTokenSecret"`
}

type Providers struct {
	Providers []Identity `yaml:"providers"`
}
