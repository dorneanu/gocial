package share

import (
	"context"
	"fmt"
	"os"

	"github.com/dorneanu/gomation/internal/entity"
)

type Service interface {
	ShareArticle(entity.ArticleShare, Repository) error
	ShareComment(entity.CommentShare, Repository) error
	GetShareRepo(entity.IdentityProvider) (Repository, error)
}

type shareService struct{}

func NewShareService() Service {
	return shareService{}
}

// ShareArticle shares an article using the specified repository
func (s shareService) ShareArticle(article entity.ArticleShare, repo Repository) error {
	// Send article to each available repository
	err := repo.ShareArticle(context.Background(), article)
	return err
}

// TODO: Implement ShareComment ...
func (s shareService) ShareComment(comment entity.CommentShare, repo Repository) error {
	return nil
}

func (s shareService) GetShareRepo(identity entity.IdentityProvider) (Repository, error) {
	if identity.Provider == "twitter" { // twitter
		twitterConfig := &TwitterConfig{
			ConsumerKey:    os.Getenv("TWITTER_CLIENT_KEY"),
			ConsumerSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
			AccessToken:    identity.AccessToken,
			AccessSecret:   identity.AccessTokenSecret,
		}
		twitterShareRepo := NewTwitterShareRepository(twitterConfig)
		return twitterShareRepo, nil

	} else if identity.Provider == "linkedin" { // linkedin
		linkedinShareRepo := NewLinkedinShareRepository(identity)
		return linkedinShareRepo, nil

	}
	return nil, fmt.Errorf("Didn't find repository")
}
