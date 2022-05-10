package entity

import "time"

// IdentityProvider holds information about an identity
// It contains sets of OAuth credentials and user specific details
// so that an user can communite to multiple (identity) providers
type IdentityProvider struct {
	Provider          string     `yaml:"provider"`
	UserName          string     `yaml:"name"`
	UserID            string     `yaml:"id"`
	AccessToken       string     `yaml:"accessToken"`
	AccessTokenSecret string     `yaml:"accessTokenSecret"`
	RefreshToken      string     `yaml:"refreshToken"`
	ExpiresAt         *time.Time `yaml:"expiry"`
}

type Providers struct {
	Providers []IdentityProvider `yaml:"providers"`
}
