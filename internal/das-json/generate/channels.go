package main

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed tmpl/channels.tmpl
var channelTmpl string

var channelTemplate = template.Must(template.New("channels").Funcs(
	template.FuncMap{
		"bt": func() string { return "`" },
	}).Parse(channelTmpl))

type Channels struct {
	Package string
}

func (s Channels) Render() ([]byte, error) {
	var buf bytes.Buffer
	if err := channelTemplate.Execute(&buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s Channels) RenderFile(path string) error {
	return RenderFile(path, s)
}
