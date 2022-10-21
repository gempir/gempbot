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
	BttvToken         string `json:"bttvToken"`
	DSN               string `json:"DSN"`
}

func FromEnv() *Config {
	return &Config{
		ClientID:     Getenv("NEXT_PUBLIC_TWITCH_CLIENT_ID"),
		ClientSecret: Getenv("TWITCH_CLIENT_SECRET"),
		Secret:       Getenv("SECRET"),
		WebBaseUrl:   Getenv("NEXT_PUBLIC_BASE_URL"),
		CookieDomain: Getenv("COOKIE_DOMAIN"),
		Username:     Getenv("TWITCH_USERNAME"),
		OAuth:        Getenv("TWITCH_OAUTH"),
		BttvToken:    Getenv("BTTV_TOKEN"),
		DSN:          Getenv("DSN"),
	}
}

func Getenv(key string) string {
	variable := os.Getenv(key)

	return strings.TrimSuffix(strings.TrimPrefix(variable, "\""), "\"")
}
