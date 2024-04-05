package providers

import (
	"context"
	"log"

	"github.com/semanser/ai-coder/assets"
	"github.com/semanser/ai-coder/config"
	"github.com/semanser/ai-coder/database"
	"github.com/semanser/ai-coder/templates"

	"github.com/tmc/langchaingo/llms"
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
	return Summary(p.client, config.Config.OpenAIModel, query, n)
}

func (p OpenAIProvider) DockerImageName(task string) (string, error) {
	return DockerImageName(p.client, config.Config.OpenAIModel, task)
}

func (p OpenAIProvider) NextTask(args NextTaskOptions) *database.Task {
	log.Println("Getting next task")

	promptArgs := map[string]interface{}{
		"DockerImage":     args.DockerImage,
		"ToolPlaceholder": "Always use your function calling functionality, instead of returning a text result.",
		"Tasks":           args.Tasks,
	}

	prompt, err := templates.Render(assets.PromptTemplates, "prompts/agent.tmpl", promptArgs)

	// TODO In case of lots of tasks, we should try to get a summary using gpt-3.5
	if len(prompt) > 30000 {
		log.Println("Prompt too long, asking user")
		return defaultAskTask("My prompt is too long and I can't process it")
	}

	if err != nil {
		log.Println("Failed to render prompt, asking user, %w", err)
		return defaultAskTask("There was an error getting the next task")
	}

	messages := tasksToMessages(args.Tasks, prompt)

	resp, err := p.client.GenerateContent(
		context.Background(),
		messages,
		llms.WithTemperature(0.0),
		llms.WithModel(p.model),
		llms.WithTopP(0.2),
		llms.WithN(1),
		llms.WithTools(Tools),
	)

	if err != nil {
		log.Printf("Failed to get response from model %v", err)
		return defaultAskTask("There was an error getting the next task")
	}

	task, err := toolToTask(resp.Choices)

	if err != nil {
		log.Printf("Failed to convert tool to task %v", err)
		return defaultAskTask("There was an error getting the next task")
	}

	return task
}
