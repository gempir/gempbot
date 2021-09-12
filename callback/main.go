package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gempir/bitraft/pkg/config"
	"github.com/gempir/bitraft/pkg/helix"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "text/plain"},
		Body:            "Hello",
		IsBase64Encoded: false,
	}, nil
}

var (
	helixClient *helix.Client
)

func main() {
	cfg := config.FromEnv()
	helix.NewClient(cfg.ClientID, cfg.ClientSecret, cfg.ApiBaseUrl+"/api/callback", cfg.Secret)
	lambda.Start(handler)
}

// type tokenClaims struct {
// 	UserID         string
// 	StandardClaims jwt.StandardClaims
// }

// func (t *tokenClaims) Valid() error {
// 	return nil
// }

// func handleCallback(r events.APIGatewayProxyRequest) error {
// 	code := r.QueryStringParameters["code"]

// 	resp, err := helixClient.Client.RequestUserAccessToken(code)
// 	if err != nil || resp.StatusCode >= 400 {
// 		log.Errorf("failed to request userAccessToken: %s %s", err, resp.ErrorMessage)
// 		return dashboardRedirect(c, "")
// 	}

// 	// validate
// 	success, validateResp, err := helixClient.Client.ValidateToken(resp.Data.AccessToken)
// 	if !success || err != nil {
// 		return fmt.Errorf("failed to veryify new Token %s", err)
// 	}

// 	token, err := createApiToken(validateResp.Data.UserID)
// 	if err != nil {
// 		return fmt.Errorf("failed to create jwt token in callback %s", err)
// 	}

// 	err = db.SaveUserAccessToken(validateResp.Data.UserID, resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " "))
// 	if err != nil {
// 		return fmt.Errorf("failed to set userAccessToken in callback: %s", err)
// 	}

// 	err = db.SaveBotConfig(store.BotConfig{OwnerTwitchID: validateResp.Data.UserID, JoinBot: true})
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	// s.store.PublishIngesterMessage(store.IngesterMsgJoin, validateResp.Data.Login)

// 	// go s.subscribePredictions(validateResp.Data.UserID)

// 	return dashboardRedirect(c, token)
// }

// func createApiToken(userID string) (string, error) {
// 	expirationTime := time.Now().Add(365 * 24 * time.Hour)
// 	claims := &tokenClaims{
// 		UserID: userID,
// 		StandardClaims: jwt.StandardClaims{
// 			// In JWT, the expiry time is expressed as unix milliseconds
// 			ExpiresAt: expirationTime.Unix(),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString([]byte(s.cfg.Secret))

// 	return tokenString, err
// }

// func dashboardRedirect(c echo.Context, scToken string) error {
// 	cookie := http.Cookie{
// 		Name:    "scToken",
// 		Value:   scToken,
// 		Domain:  cfg.CookieDomain,
// 		Expires: time.Now().Add(365 * 24 * time.Hour),
// 		Path:    "/",
// 	}

// 	c.SetCookie(&cookie)
// 	err := c.Redirect(http.StatusFound, s.cfg.WebBaseUrl)
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}

// 	return nil
// }

// func authenticate(c echo.Context) (nickHelix.ValidateTokenResponse, store.UserAccessToken, error) {
// 	scToken := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")

// 	// Initialize a new instance of `Claims`
// 	claims := &tokenClaims{}

// 	// Parse the JWT string and store the result in `claims`.
// 	// Note that we are passing the key in this method as well. This method will return an error
// 	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
// 	// or if the signature does not match
// 	tkn, err := jwt.ParseWithClaims(scToken, claims, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(s.cfg.Secret), nil
// 	})
// 	if err != nil || !tkn.Valid {
// 		log.Errorf("found to validate jwt: %s", err)
// 		return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, echo.NewHTTPError(http.StatusUnauthorized, "bad authentication")
// 	}

// 	token, err := s.db.GetUserAccessToken(claims.UserID)
// 	if err != nil {
// 		log.Errorf("Failed to get userAccessTokenData: %s", err.Error())
// 		return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, echo.NewHTTPError(http.StatusUnauthorized, "Failed to get userAccessTokenData: %s", err.Error())
// 	}

// 	success, resp, err := s.helixClient.Client.ValidateToken(token.AccessToken)
// 	if !success || err != nil {
// 		if err != nil {
// 			log.Errorf("token did not validate: %s", err)
// 		}

// 		// Token might be expired, let's try refreshing
// 		if resp.Error == "Unauthorized" {
// 			err := s.refreshToken(token)
// 			if err != nil {
// 				return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, echo.NewHTTPError(http.StatusUnauthorized, "failed to refresh token")
// 			}

// 			refreshedToken, err := s.db.GetUserAccessToken(claims.UserID)
// 			if err != nil {
// 				log.Errorf("Failed to get userAccessTokenData: %s", err.Error())
// 				return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, echo.NewHTTPError(http.StatusUnauthorized, "Failed to get userAccessTokenData: %s", err.Error())
// 			}

// 			success, resp, err = s.helixClient.Client.ValidateToken(refreshedToken.AccessToken)
// 			if !success || err != nil {
// 				if err != nil {
// 					log.Errorf("refreshed Token did not validate: %s", err)
// 				}

// 				return nickHelix.ValidateTokenResponse{}, refreshedToken, echo.NewHTTPError(http.StatusUnauthorized, "refreshed token did not validate")
// 			}

// 			return *resp, refreshedToken, nil
// 		}

// 		return nickHelix.ValidateTokenResponse{}, store.UserAccessToken{}, echo.NewHTTPError(http.StatusUnauthorized, "token not valid")
// 	}

// 	return *resp, token, nil
// }
