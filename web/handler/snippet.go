package handler

import (
	"diskhub/web/render"
	"diskhub/web/snippet"
	"html/template"
	"net/http"

	"github.com/esmyxvatu/feather"
)

type SnippetPage struct {
	Cookies        map[string]any
	SearchTerm     string
	IdSelected     string
	SnippetChoosed snippet.Snippet
	SnippetFound   []snippet.Snippet
	SnippetContent template.HTML
}

func SnippetHandler(ctx *feather.Context) {
	cookies, err := makeCookieMap(ctx)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	choosed := snippet.Get(ctx.Query("id"))

	data := SnippetPage{
		Cookies:        cookies,
		SearchTerm:     ctx.Query("q"),
		IdSelected:     ctx.Query("id"),
		SnippetChoosed: choosed,
		SnippetFound:   snippet.GetAll(),
		SnippetContent: render.StringToChroma(choosed.Content),
	}

	// Send the template
	ctx.Template([]string{"website/snippet.html"}, data, nil)
}
