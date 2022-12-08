package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/dorneanu/gocial/internal/identity"
	"github.com/dorneanu/gocial/internal/oauth"
	"github.com/dorneanu/gocial/internal/share"
	"github.com/dorneanu/gocial/server"
	"github.com/labstack/echo/v4"
)

var (
	// echoLambda *echoadapter.EchoLambdaV2
	echoLambda *echoadapter.EchoLambda
)

// func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
// 	return echoLambda.ProxyWithContext(ctx, req)
// }

func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	e := echo.New()

	webServerConf := server.HTTPServerConfig{
		ListenAddr:      "gocial.netlify.app",
		TokenSigningKey: "secret key",
		TokenExpiration: 5,
	}
	// TODO: Is this still needed
	// idRepo := entity.NewFileIdentityRepo("./auth.yaml")

	// Create OAuth configs for different providers
	oauthConfigs := []oauth.OAuthConfig{
		oauth.OAuthConfig{
			ProviderName: "linkedin",
			Scopes:       []string{"r_emailaddress", "r_liteprofile", "w_member_social"},
			ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
			ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
			CallbackURL:  fmt.Sprintf("https://%s/auth/callback/linkedin", webServerConf.ListenAddr),
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
	gothRepository := oauth.NewGothRepository(providerIndex, webServerConf.TokenSigningKey)

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

	// e.GET("/bla", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Alles klar")
	// })

	// e.GET("/kuku", func(c echo.Context) error {
	// 	return c.JSON(200, &echo.Map{"data": "Hello from Echo & mongoDB"})
	// })

	// echoLambda = echoadapter.NewV2(e)
	echoLambda = echoadapter.New(e)

}
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
