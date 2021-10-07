package config

import (
	"os"
	"strings"

	"github.com/gempir/gempbot/pkg/log"
	"github.com/sirupsen/logrus"
)

// Config application configuratin
type Config struct {
	Username          string `json:"username"`
	OAuth             string `json:"oauth"`
	ClientID          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	LogLevel          string `json:"logLevel"`
	Secret            string `json:"secret"`
	ApiBaseUrl        string `json:"apiBaseUrl"`
	WebBaseUrl        string `json:"webBaseUrl"`
	WebhookApiBaseUrl string `json:"webhookApiBaseUrl"`
	CookieDomain      string `json:"cookieDomain"`
	BttvToken         string `json:"bttvToken"`
	SevenTvToken      string `json:"sevenTvToken"`
	DbHost            string `json:"DbHost"`
	DbUsername        string `json:"DbUsername"`
	DbPassword        string `json:"DbPassword"`
	DbName            string `json:"DbName"`
	Environment       string `json:"environment"`
}

func FromEnv() *Config {
	baseUrl := Getenv("NEXT_PUBLIC_BASE_URL")
	cookieDomain := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(baseUrl, "https://"), "http://"), ":3000")

	webhookApiBaseUrl := Getenv("WEBHOOK_BASE_URL")
	if webhookApiBaseUrl == "" {
		webhookApiBaseUrl = baseUrl
	}

	return &Config{
		ClientID:          Getenv("TWITCH_CLIENT_ID"),
		ClientSecret:      Getenv("TWITCH_CLIENT_SECRET"),
		Secret:            Getenv("SECRET"),
		ApiBaseUrl:        baseUrl,
		WebBaseUrl:        baseUrl,
		WebhookApiBaseUrl: webhookApiBaseUrl,
		CookieDomain:      cookieDomain,
		DbHost:            Getenv("PLANETSCALE_DB_HOST"),
		DbUsername:        Getenv("PLANETSCALE_DB_USERNAME"),
		DbPassword:        Getenv("PLANETSCALE_DB_PASSWORD"),
		DbName:            Getenv("PLANETSCALE_DB"),
		Username:          Getenv("TWITCH_USERNAME"),
		OAuth:             Getenv("TWITCH_OAUTH"),
		BttvToken:         Getenv("BTTV_TOKEN"),
		SevenTvToken:      Getenv("SEVEN_TV_TOKEN"),
		Environment:       Getenv("VERCEL_ENV"),
	}
}

func Getenv(key string) string {
	variable := os.Getenv(key)

	return strings.TrimSuffix(strings.TrimPrefix(variable, "\""), "\"")
}

func (cfg *Config) setupLogger() {
	switch cfg.LogLevel {
	case "fatal":
		log.SetLogLevel(logrus.FatalLevel)
	case "panic":
		log.SetLogLevel(logrus.PanicLevel)
	case "error":
		log.SetLogLevel(logrus.ErrorLevel)
	case "warn":
		log.SetLogLevel(logrus.WarnLevel)
	case "info":
		log.SetLogLevel(logrus.InfoLevel)
	case "debug":
		log.SetLogLevel(logrus.DebugLevel)
	}
}
