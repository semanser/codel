package assets

import (
	"embed"
)

var PromptTemplates embed.FS
var ScriptTemplates embed.FS

func Init(promptTemplates embed.FS, scriptTemplates embed.FS) {
	PromptTemplates = promptTemplates
	ScriptTemplates = scriptTemplates
}
