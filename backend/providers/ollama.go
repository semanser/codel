package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/semanser/ai-coder/assets"
	"github.com/semanser/ai-coder/config"
	"github.com/semanser/ai-coder/database"
	"github.com/semanser/ai-coder/templates"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type OllamaProvider struct {
	client  *ollama.LLM
	model   string
	baseURL string
	name    ProviderType
}

func (p OllamaProvider) New() Provider {
	model := config.Config.OllamaModel
	baseURL := config.Config.OllamaServerURL

	client, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(baseURL),
		ollama.WithFormat("json"),
	)

	if err != nil {
		log.Fatalf("Failed to create Ollama client: %v", err)
	}

	return OllamaProvider{
		client:  client,
		model:   model,
		baseURL: baseURL,
		name:    ProviderOllama,
	}
}

func (p OllamaProvider) Name() ProviderType {
	return p.name
}

func (p OllamaProvider) Summary(query string, n int) (string, error) {
	model := config.Config.OllamaModel
	baseURL := config.Config.OllamaServerURL

	client, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(baseURL),
	)

	if err != nil {
		return "", fmt.Errorf("failed to create Ollama client: %v", err)
	}

	return Summary(client, p.model, query, n)
}

func (p OllamaProvider) DockerImageName(task string) (string, error) {
	model := config.Config.OllamaModel
	baseURL := config.Config.OllamaServerURL

	client, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(baseURL),
	)

	if err != nil {
		return "", fmt.Errorf("failed to create Ollama client: %v", err)
	}

	return DockerImageName(client, p.model, task)
}

type Call struct {
	Tool    string            `json:"tool"`
	Input   map[string]string `json:"tool_input"`
	Message string            `json:"message"`
}

func (p OllamaProvider) NextTask(args NextTaskOptions) *database.Task {
	log.Println("Getting next task")

	promptArgs := map[string]interface{}{
		"DockerImage":     args.DockerImage,
		"ToolPlaceholder": getToolPlaceholder(),
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
	)

	if err != nil {
		log.Printf("Failed to get response from model %v", err)
		return defaultAskTask("There was an error getting the next task")
	}

	choices := resp.Choices

	if len(choices) == 0 {
		log.Println("No choices found, asking user")
		return defaultAskTask("Looks like I couldn't find a task to run")
	}

	task, err := textToTask(choices[0].Content)

	if err != nil {
		log.Println("Failed to convert text to the next task, asking user")
		return defaultAskTask("There was an error getting the next task")
	}

	return task
}

func getToolPlaceholder() string {
	bs, err := json.Marshal(Tools)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(`You have access to the following tools:

%s

To use a tool, respond with a JSON object with the following structure: 
{
  "tool": <name of the called tool>,
  "tool_input": <parameters for the tool matching the above JSON schema>,
  "message": <a message that will be displayed to the user>
}

Always use a tool. Always reply with valid JOSN. Always include a message.
`, string(bs))
}
