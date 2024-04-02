package providers

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
