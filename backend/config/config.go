package config

import (
	"log"

	"github.com/caarlos0/env/v10"
)

type config struct {
	OpenAiKey   string `env:"OPEN_AI_KEY"`
	DatabaseUrl string `env:"DATABASE_URL"`
	Port        int    `env:"PORT" envDefault:"8080"`
}

var Config config

func Init() {
	if err := env.Parse(&Config); err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}
}
