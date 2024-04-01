package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type config struct {
	OpenAIKey       string `env:"OPEN_AI_KEY"`
	OpenAIModel     string `env:"OPEN_AI_MODEL" envDefault:"gpt-4-0125-preview"`
	OpenAIServerURL string `env:"OPEN_AI_SERVER_URL" envDefault:"https://api.openai.com/v1"`
	DatabaseURL     string `env:"DATABASE_URL" envDefault:"database.db"`
	Port            int    `env:"PORT" envDefault:"8080"`
}

var Config config

func Init() {
	godotenv.Load()

	if err := env.ParseWithOptions(&Config, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}
}
