package providers

import (
	"log"

	"github.com/semanser/ai-coder/config"
	"github.com/semanser/ai-coder/database"

	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAIProvider struct {
	client  *openai.LLM
	model   string
	baseURL string
	name    ProviderType
}

func (p OpenAIProvider) New() Provider {
	model := config.Config.OpenAIModel
	baseURL := config.Config.OpenAIServerURL

	client, err := openai.New(
		openai.WithToken(config.Config.OpenAIKey),
		openai.WithModel(model),
		openai.WithBaseURL(baseURL),
	)

	if err != nil {
		log.Fatalf("Failed to create OpenAI client: %v", err)
	}

	return OpenAIProvider{
		client:  client,
		model:   model,
		baseURL: baseURL,
		name:    ProviderOpenAI,
	}
}

func (p OpenAIProvider) Name() ProviderType {
	return p.name
}

func (p OpenAIProvider) Summary(query string, n int) (string, error) {
	// TODO Use more basic model for this task
	return Summary(p.client, config.Config.OpenAIModel, query, n)
}

func (p OpenAIProvider) DockerImageName(task string) (string, error) {
	// TODO Use more basic model for this task
	return DockerImageName(p.client, config.Config.OpenAIModel, task)
}

func (p OpenAIProvider) NextTask(args NextTaskOptions) *database.Task {
	return NextTask(args, p.client)
}
