package handler

import (
	"diskhub/web/logger"
	"diskhub/web/models"
	"diskhub/web/render"
	"html/template"
	"net/http"
	"os"

	"github.com/esmyxvatu/feather"
)

type ProjectPage struct {
	Name    string
	Id      string
	About   string
	Links   map[string]string
	Files   []models.FileObject
	Langs   []models.StringInt
	Tags    []string
	LOC     int
	Saved   bool
	Readme  template.HTML
	Path    string
	Status  string
	Cookies map[string]any
}

func ProjectHandler(ctx *feather.Context) {
	// Search the requested project
	var project models.Project
	found := false
	for _, proj := range models.Projects {
		if proj.Id == ctx.Params["id"] {
			project = proj
			found = true
			break
		}
	}
	if !found {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	// Load the ReadMe content from the file
	content, err := os.ReadFile(project.ReadMe)
	if err != nil {
		logger.Console.Error("Error reading file: %s", err.Error())
		content = []byte{42, 42, 78, 111, 32, 82, 101, 97, 100, 109, 101, 32, 70, 111, 117, 110, 100, 42, 42}
	}

	// Some code for the template
	project.Langs[len(project.Langs)-1].IsLast = true

	cookies, err := makeCookieMap(ctx)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	// Unify the data
	data := ProjectPage{
		Name:    project.Name,
		Id:      project.Id,
		About:   project.About,
		Links:   project.Links,
		Files:   project.Files,
		Langs:   project.Langs,
		LOC:     project.LOC,
		Saved:   project.Saved,
		Tags:    project.Tags,
		Readme:  render.MarkdownToHTML(content),
		Path:    project.Path,
		Status:  project.Status,
		Cookies: cookies,
	}

	// Send the data
	ctx.Template([]string{"website/project.html"}, data, nil)
}
