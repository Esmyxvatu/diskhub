package render

import (
	"bytes"
	"diskhub/web/logger"
	"html/template"

	"github.com/yuin/goldmark"
	mdHighlight "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"go.abhg.dev/goldmark/wikilink"
)

type customResolver struct{}

func (customResolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	var _hash = []byte("#")
	dest := make([]byte, len(n.Target)+len(_hash)+len(n.Fragment))

	var i int
	if len(n.Target) > 0 {
		i += copy(dest, n.Target)
	}

	if len(n.Fragment) > 0 {
		i += copy(dest[i:], _hash)
		i += copy(dest[i:], n.Fragment)
	}

	return dest[:i], nil
}

func MarkdownToHTML(mdContent []byte) template.HTML {
	// Convert markdown to HTML
	// htmlContent := blackfriday.Run(mdContent)

	myresolver := customResolver{}
	markdown := goldmark.New(
		goldmark.WithExtensions(
			mdHighlight.Highlighting,
			extension.Table,
			&wikilink.Extender{Resolver: myresolver},
		),
	)
	var byteContent bytes.Buffer
	if err := markdown.Convert(mdContent, &byteContent); err != nil {
		logger.Console.Fatal("%v", err)
	}
	htmlContent := byteContent.String()

	return template.HTML(htmlContent)
}
