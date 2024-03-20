package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/invopop/jsonschema"
	openai "github.com/sashabaranov/go-openai"
	"github.com/semanser/ai-coder/assets"
	"github.com/semanser/ai-coder/templates"
)

type CommandType string

const (
	Input    CommandType = "input"
	Terminal CommandType = "terminal"
	Browser  CommandType = "browser"
	Code     CommandType = "code"
	Ask      CommandType = "ask"
	Done     CommandType = "done"
)

type Command struct {
	ID          int         `json:"id"`
	Type        CommandType `json:"type"`
	Args        interface{} `json:"args,omitempty"`
	Result      interface{} `json:"result,omitempty"`
	Description string      `json:"description"`
}

var openAIclient *openai.Client
var OPEN_AI_KEY string

func Init() {
	OPEN_AI_KEY := os.Getenv("OPEN_AI_KEY")
	openAIclient = openai.NewClient(OPEN_AI_KEY)

	if OPEN_AI_KEY == "" {
		log.Fatal("OPEN_AI_KEY is not set")
	}

	commands := []Command{
		{
			ID:   1,
			Type: Input,
			Args: InputArgs{
				Query: "Create a new game of tic-tac-toe in react",
			},
		},
		{
			ID:   2,
			Type: Terminal,
			Args: TerminalArgs{
				Input:       "npx create-react-app tic-tac-toe",
				Description: "I'm trying to create a new react app using create-react-app template.",
			},
			Result: "Create React App was successfully created!",
		},
	}

	c, err := NextCommand(AgentPrompt{
		Commands: commands,
	})

	if err != nil {
		log.Fatalf("Failed to get next command: %v", err)
	}

	log.Printf("Command: %v", c)
	log.Printf("Command Args: %v", c.Args)
}

type Description string

type InputArgs struct {
	Query string
}

type TerminalArgs struct {
	Input string
	Description
}

type BrowserAction string

const (
	Read BrowserAction = "read"
	Url  BrowserAction = "url"
)

type BrowserArgs struct {
	Url    string
	Action BrowserAction
	Description
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
	Description
}

type AskArgs struct {
	Input string
	Description
}

type DoneArgs struct {
	Description
}

type AgentPrompt struct {
	Commands []Command
}

func NextCommand(args AgentPrompt) (*Command, error) {
	prompt, err := templates.Render(assets.PromptTemplates, "prompts/agent.tmpl", args)
	if err != nil {
		return nil, err
	}

	tools := []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        string(Terminal),
				Description: "Calls a terminal command",
				Parameters:  jsonschema.Reflect(&TerminalArgs{}).Definitions["TerminalArgs"],
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        string(Browser),
				Description: "Opens a browser to loop for additional information",
				Parameters:  jsonschema.Reflect(&BrowserArgs{}).Definitions["BrowserArgs"],
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        string(Code),
				Description: "Modifies or retrieves code files",
				Parameters:  jsonschema.Reflect(&CodeArgs{}).Definitions["CodeArgs"],
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        string(Ask),
				Description: "Sends a question to the user for additional information",
				Parameters:  jsonschema.Reflect(&AskArgs{}).Definitions["AskArgs"],
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        string(Done),
				Description: "Mark the whole task as done. Should be called at the very end when everything is completed",
				Parameters:  jsonschema.Reflect(&DoneArgs{}).Definitions["DoneArgs"],
			},
		},
	}

	req := openai.ChatCompletionRequest{
		Temperature: 0.0,
		Model:       openai.GPT3Dot5Turbo0125,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Tools: tools,
	}

	resp, err := openAIclient.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("completion error: %v", err)
	}

	tool := resp.Choices[0].Message.ToolCalls[0]

	if tool.Function.Name == "" {
		return nil, fmt.Errorf("no tool found")
	}

	command := Command{
		Type: CommandType(tool.Function.Name),
	}

	switch tool.Function.Name {
	case string(Terminal):
		params, err := extractArgs(tool.Function.Arguments, &TerminalArgs{})
		if err != nil {
			return nil, fmt.Errorf("failed to extract terminal args: %v", err)
		}
		command.Args = params.Input
		command.Description = string(params.Description)

	case string(Browser):
		params, err := extractArgs(tool.Function.Arguments, &BrowserArgs{})
		if err != nil {
			return nil, fmt.Errorf("failed to extract browser args: %v", err)
		}
		command.Args = params
		command.Description = string(params.Description)
	case string(Code):
		params, err := extractArgs(tool.Function.Arguments, &CodeArgs{})
		if err != nil {
			return nil, fmt.Errorf("failed to extract code args: %v", err)
		}
		command.Args = params
		command.Description = string(params.Description)
	case string(Ask):
		params, err := extractArgs(tool.Function.Arguments, &AskArgs{})
		if err != nil {
			return nil, fmt.Errorf("failed to extract ask args: %v", err)
		}
		command.Args = params
		command.Description = string(params.Description)
	case string(Done):
		params, err := extractArgs(tool.Function.Arguments, &DoneArgs{})
		if err != nil {
			return nil, fmt.Errorf("failed to extract ask args: %v", err)
		}
		command.Args = params
		command.Description = string(params.Description)
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
