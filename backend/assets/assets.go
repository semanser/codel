package assets

import (
	"embed"
)

var PromptTemplates embed.FS

func Init(promptTemplates embed.FS) {
	PromptTemplates = promptTemplates
}
