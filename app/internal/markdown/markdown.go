package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/FurqanSoftware/goldmark-katex"
)

var engine = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Table,
		extension.Strikethrough,
		extension.Linkify,
		extension.TaskList,
		extension.Footnote,
		&katex.Extender{},
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
	),
)

func Render(src string) (string, error) {
	var buf bytes.Buffer
	if err := engine.Convert([]byte(src), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
