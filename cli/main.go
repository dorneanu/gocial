package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dorneanu/gomation/internal/auth"
	"github.com/dorneanu/gomation/internal/config"
	"github.com/dorneanu/gomation/internal/entity"
	"github.com/dorneanu/gomation/internal/share"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

var (
	postURL     string
	postTitle   string
	postComment string
	conf        config.Config
)

func main() {
	app := &cli.App{
		// Flags: globalFlags,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Victor Dorneanu",
				Email: "some e-mail",
			},
		},
		Version:  "v0.1",
		Compiled: time.Now(),
		Commands: []*cli.Command{
			{
				// authenticate sub-command
				Name:    "authenticate",
				Aliases: []string{"a"},
				Usage:   "Authenticate against identity providers",
				Action: func(c *cli.Context) error {

					// New linkedin repository
					oauthConfigs := []auth.OAuthConfig{
						auth.OAuthConfig{
							ProviderName: "linkedin",
							Scopes:       []string{"r_emailaddress", "r_liteprofile", "w_member_social"},
							ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
							ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
							CallbackURL:  "http://localhost:3000/auth/linkedin/callback",
						},
						auth.OAuthConfig{
							ProviderName: "twitter",
							ClientID:     os.Getenv("TWITTER_CLIENT_KEY"),
							ClientSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
							CallbackURL:  "http://127.0.0.1:3000/auth/twitter/callback",
						},
					}

					// New identity repository
					idRepo := entity.NewFileIdentityRepo("./auth.yaml")

					// New goth auth repository
					r := mux.NewRouter()
					gothRepository := auth.NewGothRepository(r, oauthConfigs, idRepo)

					// New auth service
					authService := auth.NewService(gothRepository, r)
					authService.Start()
					return nil
				},
			},
			{
				// post sub-command
				Name:    "post",
				Aliases: []string{"p"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "url",
						Usage:       "URL",
						Destination: &postURL,
					},
					&cli.StringFlag{
						Name:        "title",
						Usage:       "Post title",
						Destination: &postTitle,
					},
					&cli.StringFlag{
						Name:        "comment",
						Usage:       "Post commentary",
						Destination: &postComment,
					},
				},
				Usage: "Post some article",
				Action: func(c *cli.Context) error {
					authConf := entity.NewFileIdentityRepo("./auth.yaml")
					err := authConf.Load()
					if err != nil {
						return fmt.Errorf("Couldn't load auth details: %s", err)
					}

					article := entity.ArticleShare{
						URL:     postURL,
						Title:   postTitle,
						Comment: postComment,
					}

					id, err := authConf.GetByProvider("twitter")
					if err != nil {
						return fmt.Errorf("Couldn't get identity: %s", err)
					}

					// linkedShareRepo := share.NewLinkedinShareRepository(id)
					// oauth2 configures a client that uses app credentials to keep a fresh token
					twitterConfig := &share.TwitterConfig{
						ConsumerKey:    os.Getenv("TWITTER_CLIENT_KEY"),
						ConsumerSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
						AccessToken:    id.AccessToken,
						AccessSecret:   id.AccessTokenSecret,
					}
					twitterShareRepo := share.NewTwitterShareRepository(twitterConfig)

					// New share service
					shareService := share.NewService(twitterShareRepo)
					err = shareService.ShareArticle(article)
					if err != nil {
						return fmt.Errorf("Couldn't share article: %s", err)
					}
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
