package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/dorneanu/gomation/internal/entity"
	"github.com/dorneanu/gomation/internal/identity"
	"github.com/dorneanu/gomation/internal/oauth"
	"github.com/dorneanu/gomation/internal/share"
	"github.com/dorneanu/gomation/server"
	"github.com/labstack/echo/v4"
)

var (
	echoLambda *echoadapter.EchoLambdaV2
)

func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("echo cold start")
	e := echo.New()

	webServerConf := server.HTTPServerConfig{
		ListenAddr:      "127.0.0.1:3000",
		TokenSigningKey: "secret key",
		TokenExpiration: 5,
	}
	// TODO: Is this still needed
	idRepo := entity.NewFileIdentityRepo("./auth.yaml")

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
	httpServer := server.NewHTTPService(webServerConf)
	httpServer.Start(e)
	echoLambda = echoadapter.NewV2(e)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}
func main() {
	lambda.Start(Handler)
}
