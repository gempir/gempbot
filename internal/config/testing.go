package config

func NewMockConfig() *Config {
	return &Config{
		ClientID:          "clientId",
		ClientSecret:      "clientSecret",
		Secret:            "secret",
		WebBaseUrl:        "https://web.test.gempir.com",
		WebhookApiBaseUrl: "https://webhook.test.gempir.com",
		CookieDomain:      "https://test.gempir.com",
		Username:          "username",
		OAuth:             "oauth",
		BttvToken:         "bttvToken",
	}
}
