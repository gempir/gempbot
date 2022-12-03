package config

import (
	"os"
	"strings"
)

// Config application configuratin
type Config struct {
	Username          string `json:"username"`
	OAuth             string `json:"oauth"`
	ClientID          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	Secret            string `json:"secret"`
	WebBaseUrl        string `json:"webBaseUrl"`
	WebhookApiBaseUrl string `json:"webhookApiBaseUrl"`
	CookieDomain      string `json:"cookieDomain"`
	DSN               string `json:"DSN"`
	ListenAddress     string `json:"listenAddress"`
}

func FromEnv() *Config {
	webhookApiBaseUrl := Getenv("WEBHOOK_API_BASE_URL")
	if webhookApiBaseUrl == "" {
		webhookApiBaseUrl = Getenv("NEXT_PUBLIC_API_BASE_URL")
	}

	listenAddress := Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = ":3010"
	}

	return &Config{
		ClientID:          Getenv("NEXT_PUBLIC_TWITCH_CLIENT_ID"),
		ClientSecret:      Getenv("TWITCH_CLIENT_SECRET"),
		Secret:            Getenv("SECRET"),
		WebBaseUrl:        Getenv("NEXT_PUBLIC_BASE_URL"),
		WebhookApiBaseUrl: webhookApiBaseUrl,
		CookieDomain:      Getenv("COOKIE_DOMAIN"),
		Username:          Getenv("TWITCH_USERNAME"),
		OAuth:             Getenv("TWITCH_OAUTH"),
		DSN:               Getenv("DSN"),
		ListenAddress:     listenAddress,
	}
}

func Getenv(key string) string {
	variable := os.Getenv(key)

	return strings.TrimSuffix(strings.TrimPrefix(variable, "\""), "\"")
}
