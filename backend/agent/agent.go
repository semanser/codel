package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/invopop/jsonschema"
	openai "github.com/sashabaranov/go-openai"
	"github.com/semanser/ai-coder/assets"
	"github.com/semanser/ai-coder/models"
	"github.com/semanser/ai-coder/services"
	"github.com/semanser/ai-coder/templates"
)

type Message string

type InputArgs struct {
	Query string
}

type TerminalArgs struct {
	Input string
	Message
}

type BrowserAction string

const (
	Read BrowserAction = "read"
	Url  BrowserAction = "url"
)

type BrowserArgs struct {
	Url    string
	Action BrowserAction
	Input  string
	Message
}

type CodeAction string

const (
	ReadFile   CodeAction = "read_file"
	UpdateFile CodeAction = "update_file"
)

type CodeArgs struct {
	Action  CodeAction
	Content string
	Path    string
	Message
}

type AskArgs struct {
	Message
}

type DoneArgs struct {
	Message
}

type AgentPrompt struct {
	Tasks       []models.Task
	DockerImage string
}

func NextTask(args AgentPrompt) (*models.Task, error) {
	log.Println("Getting next command")

	prompt, err := templates.Render(assets.PromptTemplates, "prompts/agent.tmpl", args)

	if len(args.Tasks) > 30000 {
		return nil, fmt.Errorf("too big prompt")
	}

	if err != nil {
		return nil, err
	}

	tools := []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        string(models.Terminal),
				Description: "Calls a terminal command",
				Parameters:  jsonschema.Reflect(&TerminalArgs{}).Definitions["TerminalArgs"],
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        string(models.Browser),
				Description: "Opens a browser to look for additional information",
				Parameters:  jsonschema.Reflect(&BrowserArgs{}).Definitions["BrowserArgs"],
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        string(models.Code),
				Description: "Modifies or reads code files",
				Parameters:  jsonschema.Reflect(&CodeArgs{}).Definitions["CodeArgs"],
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        string(models.Ask),
				Description: "Sends a question to the user for additional information",
				Parameters:  jsonschema.Reflect(&AskArgs{}).Definitions["AskArgs"],
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        string(models.Done),
				Description: "Mark the whole task as done. Should be called at the very end when everything is completed",
				Parameters:  jsonschema.Reflect(&DoneArgs{}).Definitions["DoneArgs"],
			},
		},
	}

	req := openai.ChatCompletionRequest{
		Temperature: 0.0,
		Model:       openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		TopP:  0.2,
		Tools: tools,
		N:     1,
	}

	resp, err := services.OpenAIclient.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("completion error: %v", err)
	}

	choices := resp.Choices

	if len(choices) == 0 {
		return nil, fmt.Errorf("no choices found")
	}

	toolCalls := choices[0].Message.ToolCalls

	if len(toolCalls) == 0 {
		return nil, fmt.Errorf("no tool calls found")
	}

	tool := toolCalls[0]

	if tool.Function.Name == "" {
		return nil, fmt.Errorf("no tool found")
	}

	command := models.Task{
		Type: models.TaskType(tool.Function.Name),
	}

	switch tool.Function.Name {
	case string(models.Terminal):
		params, err := extractArgs(tool.Function.Arguments, &TerminalArgs{})
		if err != nil {
			return nil, fmt.Errorf("failed to extract terminal args: %v", err)
		}
		args, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal terminal args: %v", err)
		}
		command.Args = args

		// Sometimes the model returns an empty string for the message
		msg := string(params.Message)
		if msg != "" {
			msg = params.Input
		}

		command.Message = msg
		command.Status = models.InProgress

	case string(models.Browser):
		params, err := extractArgs(tool.Function.Arguments, &BrowserArgs{})
		if err != nil {
			return nil, fmt.Errorf("failed to extract browser args: %v", err)
		}
		args, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal browser args: %v", err)
		}
		command.Args = args
		command.Message = string(params.Message)
	case string(models.Code):
		params, err := extractArgs(tool.Function.Arguments, &CodeArgs{})
		if err != nil {
			return nil, fmt.Errorf("failed to extract code args: %v", err)
		}
		args, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal code args: %v", err)
		}
		command.Args = args
		command.Message = string(params.Message)
	case string(models.Ask):
		params, err := extractArgs(tool.Function.Arguments, &AskArgs{})
		if err != nil {
			return nil, fmt.Errorf("failed to extract ask args: %v", err)
		}
		args, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal ask args: %v", err)
		}
		command.Args = args
		command.Message = string(params.Message)
	case string(models.Done):
		params, err := extractArgs(tool.Function.Arguments, &DoneArgs{})
		if err != nil {
			return nil, fmt.Errorf("failed to extract ask args: %v", err)
		}
		args, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal ask args: %v", err)
		}
		command.Args = args
		command.Message = string(params.Message)
	}

	return &command, nil
}

func extractArgs[T any](openAIargs string, args *T) (*T, error) {
	err := json.Unmarshal([]byte(openAIargs), args)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal args: %v", err)
	}
	return args, nil
}
