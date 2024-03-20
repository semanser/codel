package templates

import (
	"bytes"
	"io/fs"
	"path"
	"text/template"
)

var RootFolder = "templates"

func Render(fs fs.ReadFileFS, name string, params any) (string, error) {
	p := path.Join(RootFolder, string(name))
	promptBytes, err := fs.ReadFile(p)

	if err != nil {
		return "", err
	}

	prompt := string(promptBytes)

	t := template.Must(template.New(string(name)).Parse(prompt))
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, params)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
