package config

func Default() Config {
	return Config{
		API: API{Port: 8080},
		MSTeams: MSTeamsConfig{
			URL:         "https://snappcab.webhook.office.com/webhookb2/653c55bc-7686-4829-a582-56ab1215522a@17d2e12c-c498-4570-85de-a88e58c5bb02/IncomingWebhook/bb562f28b55d44c686de7e2a01bc8ac3/b3acfd2c-e3fc-4311-8688-d0ec6420bc56/V2A0pTeg9SptTqwgaQXzTj25j2EK6xYAdbcawcN6mj8081",
			RuntimeURL:  "https://...",
			PlatformURL: "https://...",
			NetworkURL:  "https://...",
		},
	}
}
