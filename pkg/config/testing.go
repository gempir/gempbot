package config

func NewTestConfig() *Config {
	return &Config{
		ClientID:          "clientId",
		ClientSecret:      "clientSecret",
		Secret:            "secret",
		ApiBaseUrl:        "https://api.test.gempir.com",
		WebBaseUrl:        "https://web.test.gempir.com",
		WebhookApiBaseUrl: "https://webhook.test.gempir.com",
		CookieDomain:      "https://test.gempir.com",
		DbHost:            "1.1.1.1",
		DbUsername:        "dbUsername",
		DbPassword:        "dbPassword",
		DbName:            "dbName",
		Username:          "username",
		OAuth:             "oauth",
		BttvToken:         "bttvToken",
		SevenTvToken:      "sevenTvToken",
		Environment:       "environment",
		NewrelicLicense:   "newrelicLicense",
	}
}
