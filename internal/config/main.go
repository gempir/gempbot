package config

import (
	"os"
	"strings"
)

// Config application configuratin
type Config struct {
	BotUserID         string `json:"botUserId"`
	OAuth             string `json:"oauth"`
	ClientID          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	Secret            string `json:"secret"`
	WebBaseUrl        string `json:"webBaseUrl"`
	WebhookApiBaseUrl string `json:"webhookApiBaseUrl"`
	ApiBaseUrl        string `json:"apiBaseUrl"`
	CookieDomain      string `json:"cookieDomain"`
	DSN               string `json:"DSN"`
	LogLevel          string `json:"logLevel"`
	ListenAddress     string `json:"listenAddress"`
	YsweetUrl         string `json:"ysweetUrl"`
	YsweetToken       string `json:"ysweetToken"`
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

	apiBaseUrl := "http://" + listenAddress
	if strings.Contains(webhookApiBaseUrl, "gempir.com") {
		apiBaseUrl = webhookApiBaseUrl
	}

	logLevel := Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		ClientID:          Getenv("NEXT_PUBLIC_TWITCH_CLIENT_ID"),
		ClientSecret:      Getenv("TWITCH_CLIENT_SECRET"),
		Secret:            Getenv("SECRET"),
		WebBaseUrl:        Getenv("NEXT_PUBLIC_BASE_URL"),
		WebhookApiBaseUrl: webhookApiBaseUrl,
		ApiBaseUrl:        apiBaseUrl,
		CookieDomain:      Getenv("COOKIE_DOMAIN"),
		BotUserID:         Getenv("BOT_USER_ID"),
		DSN:               Getenv("DSN"),
		LogLevel:          logLevel,
		ListenAddress:     listenAddress,
		YsweetUrl:         Getenv("YSWEET_URL"),
		YsweetToken:       Getenv("YSWEET_TOKEN"),
	}
}

func Getenv(key string) string {
	variable := os.Getenv(key)

	return strings.TrimSuffix(strings.TrimPrefix(variable, "\""), "\"")
}
