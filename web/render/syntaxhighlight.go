package render

import (
	"bytes"
	"diskhub/web/logger"
	"html/template"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func StringToChroma(content string) template.HTML {
	// Define the lexer
	lexer := lexers.Analyse(content)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Load a theme style
	style := styles.Get("swapoff")
	if style == nil {
		style = styles.Fallback
	}

	// Define the formater to uses classes
	formatter := html.New(
		html.WithClasses(true),
		html.WithAllClasses(true),
	)

	// Tokenize the content
	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		logger.Console.Fatal("%s", err.Error())
	}

	// Format the tokens
	var w bytes.Buffer
	err = formatter.Format(&w, style, iterator)
	if err != nil {
		logger.Console.Fatal("%s", err.Error())
	}

	// Convert to HTML
	return template.HTML(w.String())
}
