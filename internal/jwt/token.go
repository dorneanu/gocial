package jwt

import (
	"fmt"
	"time"

	"github.com/dorneanu/gocial/internal/entity"
	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	UserName          string
	UserID            string
	UserDescription   string
	UserAvatarURL     string
	Provider          string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	jwt.StandardClaims
}

type Token struct {
	token        *jwt.Token
	signedString string
	claims       *JwtCustomClaims
}

// NewToken returns a signed JWT token
func NewToken(id entity.IdentityProvider, signingKey string) (string, error) {
	// Create the Claims
	claims := &JwtCustomClaims{
		UserName:          id.UserName,
		UserID:            id.UserID,
		UserDescription:   id.UserDescription,
		Provider:          id.Provider,
		UserAvatarURL:     id.UserAvatarURL,
		AccessToken:       id.AccessToken,
		AccessTokenSecret: id.AccessTokenSecret,
		RefreshToken:      id.RefreshToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer:    id.Provider,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed string
	ss, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("Could not sign token")
	}

	return ss, nil
}
