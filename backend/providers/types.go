package providers

type Message string

type Messanger interface {
	GetMessage() Message
}

type InputArgs struct {
	Query string
	Message
}

func (i *InputArgs) GetMessage() Message {
	return i.Message
}

type TerminalArgs struct {
	Input string
	Message
}

func (t *TerminalArgs) GetMessage() Message {
	return t.Message
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

func (b *BrowserArgs) GetMessage() Message {
	return b.Message
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

func (c *CodeArgs) GetMessage() Message {
	return c.Message
}

type AskArgs struct {
	Message
}

func (a *AskArgs) GetMessage() Message {
	return a.Message
}

type DoneArgs struct {
	Message
}

func (d *DoneArgs) GetMessage() Message {
	return d.Message
}
