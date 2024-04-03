package providers

import (
	"context"

	"github.com/semanser/ai-coder/assets"
	"github.com/semanser/ai-coder/templates"
	"github.com/tmc/langchaingo/llms"
)

func Summary(llm llms.Model, model string, query string, n int) (string, error) {
	prompt, err := templates.Render(assets.PromptTemplates, "prompts/summary.tmpl", map[string]any{
		"Text": query,
		"N":    n,
	})
	if err != nil {
		return "", err
	}

	response, err := llms.GenerateFromSinglePrompt(
		context.Background(),
		llm,
		prompt,
		llms.WithTemperature(0.0),
		// TODO Use a simpler model for this task
		llms.WithModel(model),
		llms.WithTopP(0.2),
		llms.WithN(1),
	)

	return response, err
}

func DockerImageName(llm llms.Model, model string, task string) (string, error) {
	prompt, err := templates.Render(assets.PromptTemplates, "prompts/docker.tmpl", map[string]any{
		"Task": task,
	})
	if err != nil {
		return "", err
	}

	response, err := llms.GenerateFromSinglePrompt(
		context.Background(),
		llm,
		prompt,
		llms.WithTemperature(0.0),
		llms.WithModel(model),
		llms.WithTopP(0.2),
		llms.WithN(1),
	)

	return response, err
}
