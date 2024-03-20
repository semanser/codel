package agent

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/semanser/ai-coder/assets"
	"github.com/semanser/ai-coder/templates"
)

var openAIclient *openai.Client
var OPEN_AI_KEY string

func Init() {
	OPEN_AI_KEY := os.Getenv("OPEN_AI_KEY")
	openAIclient = openai.NewClient(OPEN_AI_KEY)

	if OPEN_AI_KEY == "" {
		log.Fatal("OPEN_AI_KEY is not set")
	}
}

func Request() {
	prompt, err := templates.Render(assets.PromptTemplates, "prompts/agent.tmpl", "")
	if err != nil {
		log.Printf("Error rendering prompt: %v\n", err)
		return
	}

	req := openai.ChatCompletionRequest{
		Temperature: 0.0,
		Model:       openai.GPT3Dot5Turbo16K0613,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		// Functions: []openai.FunctionDefinition{
		// 	{
		// 		Name:        "ParseDataToJSON",
		// 		Description: "Parses text data from the webpage to JSON format",
		// 		Parameters:  requestBody.Schema,
		// 	},
		// },
	}

	resp, err := openAIclient.CreateChatCompletion(context.Background(), req)
	if err != nil {
		log.Printf("Completion error: %v\n", err)
		return
	}

	response := resp.Choices[0].Message.Content

	log.Printf("Response: %v\n", response)
}
