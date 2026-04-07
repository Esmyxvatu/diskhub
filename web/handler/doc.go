package handler

import (
	"diskhub/web/logger"
	"diskhub/web/models"
	"html/template"
	"net/http"
	"strings"

	"github.com/esmyxvatu/feather"
)

type DocPage struct {
	Articles []models.DocElement
	Origin   string
	Project  models.Project
	Name     string
	Content  template.HTML
	Cookies  map[string]any
}

func DocShowerHandler(ctx *feather.Context) {
	// Search for the project
	var projecte models.Project
	found := false
	for _, proj := range models.Projects {
		if proj.Id == ctx.Params["id"] {
			projecte = proj
			found = true
			break
		}
	}
	if !found {
		logger.Console.Error("Project not found for ID: %s", ctx.Params["id"])
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	// Get the list of all articles of the doc and the content of the one asked
	contentMd := template.HTML("<pre> Unable to find the content of article </pre>")
	articleAsked := projecte.Wiki.PathMap[strings.ToLower(ctx.Params["page"])]
	if articleAsked != nil {
		contentMd = articleAsked.Content
	}

	cookies, err := makeCookieMap(ctx)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	name := ctx.Params["page"]
	parts := strings.Split(name, "/")
	name = strings.Join(parts, " > ")

	// Unify all the data
	data := DocPage{
		Project:  projecte,
		Origin:   name,
		Articles: projecte.Wiki.Articles,
		Name:     ctx.Params["page"],
		Content:  contentMd,
		Cookies:  cookies,
	}

	// Send the template
	ctx.Template([]string{"website/docs.html"}, data, nil)
}
