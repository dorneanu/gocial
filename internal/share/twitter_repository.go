package share

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dorneanu/gomation/internal/entity"
)

// https://developer.twitter.com/en/docs/counting-characters
const twitterMaxCharacters = 280

// TwitterShareRepository implements share.Repository
type TwitterShareRepository struct {
	client *twitter.Client
}

type TwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func NewTwitterShareRepository(twitterConf *TwitterConfig) *TwitterShareRepository {
	// Create new twitter client based on the oauth config
	//
	// https://developer.twitter.com/en/docs/authentication/oauth-1-0a
	config := oauth1.NewConfig(twitterConf.ConsumerKey, twitterConf.ConsumerSecret)
	token := oauth1.NewToken(twitterConf.AccessToken, twitterConf.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	return &TwitterShareRepository{
		client: client,
	}
}

// ShareArticle sends a new Tweet
func (t *TwitterShareRepository) ShareArticle(ctx context.Context, article entity.ArticleShare) error {
	// Compose post
	post := fmt.Sprintf("%s %s", article.Comment, article.URL)

	// Check post length
	if len(post) > twitterMaxCharacters {
		return fmt.Errorf("Post max characters exceeded: %d (allowed: %d)", len(post), twitterMaxCharacters)
	}

	// Send a Tweet
	_, _, err := t.client.Statuses.Update(post, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	return nil
}
