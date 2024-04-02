package providers

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/semanser/ai-coder/database"
)

type ProviderType string

const (
	ProviderOpenAI ProviderType = "openai"
)

type Provider interface {
	New() Provider
	Summary(query string, n int) (string, error)
	DockerImageName(task string) (string, error)
	NextTask(args NextTaskOptions) *database.Task
}

type NextTaskOptions struct {
	Tasks       []database.Task
	DockerImage string
}

func ProviderFactory(provider ProviderType) (Provider, error) {
	switch provider {
	case "openai":
		return OpenAIProvider{}.New(), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

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

func defaultAskTask(message string) *database.Task {
	task := database.Task{
		Type: database.StringToNullString("ask"),
	}

	task.Args = database.StringToNullString("{}")
	task.Message = sql.NullString{
		String: fmt.Sprintf("%s. What should I do next?", message),
		Valid:  true,
	}

	return &task
}

func extractArgs[T any](openAIargs string, args *T) (*T, error) {
	err := json.Unmarshal([]byte(openAIargs), args)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal args: %v", err)
	}
	return args, nil
}
