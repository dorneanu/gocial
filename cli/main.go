package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dorneanu/gocial/internal/entity"
	"github.com/dorneanu/gocial/internal/identity"
	"github.com/dorneanu/gocial/internal/oauth"
	"github.com/dorneanu/gocial/internal/share"
	"github.com/dorneanu/gocial/server"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
)

var (
	postURL     string
	postTitle   string
	postComment string
)

func main() {
	app := &cli.App{
		// Flags: globalFlags,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Victor Dorneanu",
				Email: "info ät dornea DOT nu",
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
					webServerConf := server.HTTPServerConfig{
						ListenAddr:      "127.0.0.1:3000",
						TokenSigningKey: "secret key",
						TokenExpiration: 5,
					}

					// Create OAuth configs for different providers
					oauthConfigs := []oauth.OAuthConfig{
						oauth.OAuthConfig{
							ProviderName: "linkedin",
							Scopes:       []string{"r_emailaddress", "r_liteprofile", "w_member_social"},
							ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
							ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
							CallbackURL:  fmt.Sprintf("http://%s/auth/callback/linkedin", webServerConf.ListenAddr),
						},
						oauth.OAuthConfig{
							ProviderName: "twitter",
							ClientID:     os.Getenv("TWITTER_CLIENT_KEY"),
							ClientSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
							CallbackURL:  fmt.Sprintf("http://%s/auth/callback/twitter", webServerConf.ListenAddr),
						},
					}

					// New identity repository
					// TODO: Is this still needed
					idRepo := entity.NewFileIdentityRepo("./auth.yaml")

					// New goth auth repository
					providerIndex := oauth.SetupAuthProviders(oauthConfigs)
					gothRepository := oauth.NewGothRepository(providerIndex, idRepo, webServerConf.TokenSigningKey)

					// New identity repository
					cookieIdentityRepo := identity.NewCookieIdentityRepository(&identity.CookieIdentityOptions{
						BaseCookieName:  "gocial",
						TokenSigningKey: webServerConf.TokenSigningKey,
					})

					// New OAuth authentication service service
					oauthService := oauth.NewService(
						oauth.ServiceConfig{
							Repo:          gothRepository,
							ProviderIndex: providerIndex,
						},
					)

					// New share service
					webServerConf.OAuthService = oauthService
					webServerConf.IdentityService = cookieIdentityRepo
					webServerConf.ProviderIndex = &providerIndex
					webServerConf.ShareService = share.NewShareService()

					// New web server
					e := echo.New()
					httpServer := server.NewHTTPService(webServerConf)
					httpServer.Start(e)

					// Start server
					e.Logger.Fatal(e.Start(webServerConf.ListenAddr))

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
					// authConf := entity.NewFileIdentityRepo("./auth.yaml")
					// err := authConf.Load()
					// if err != nil {
					// 	return fmt.Errorf("Couldn't load auth details: %s", err)
					// }

					// article := entity.ArticleShare{
					// 	URL:     postURL,
					// 	Title:   postTitle,
					// 	Comment: postComment,
					// }

					// id, err := authConf.GetByProvider("twitter")
					// if err != nil {
					// 	return fmt.Errorf("Couldn't get identity: %s", err)
					// }

					// linkedShareRepo := share.NewLinkedinShareRepository(id)
					// oauth2 configures a client that uses app credentials to keep a fresh token
					// twitterConfig := &share.TwitterConfig{
					// 	ConsumerKey:    os.Getenv("TWITTER_CLIENT_KEY"),
					// 	ConsumerSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
					// 	AccessToken:    id.AccessToken,
					// 	AccessSecret:   id.AccessTokenSecret,
					// }
					// twitterShareRepo := share.NewTwitterShareRepository(twitterConfig)

					// New share service
					// shareService := share.NewShareService(twitterShareRepo)
					// err = shareService.ShareArticle(article)
					// if err != nil {
					// 	return fmt.Errorf("Couldn't share article: %s", err)
					// }
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
