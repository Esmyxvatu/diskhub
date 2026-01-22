package handler

import (
	"diskhub/web/language"
	"diskhub/web/models"
	"diskhub/web/render"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/esmyxvatu/feather"
)

type SingleFilePage struct {
	Name     string
	Id       string
	Files    []models.FileObject
	Path     string
	Content  template.HTML
	FileName string
	Info     map[string]any
	Cookies  map[string]any
}

func FileShowerHandler(ctx *feather.Context) {
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
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	// Search for the file
	var file models.FileObject
	found = false
	parts := strings.Split(ctx.Params["filepath"], "/") // Split the path | "/abc/114/e.html" => [abc 114 e.html]
	if len(parts) > 1 {
		files := projecte.Files

		for _, part := range parts { // Search for every parts
			for _, fily := range files { // Search for every file
				if fily.Name == part { // Try if it's the good file
					if fily.IsDir { // If it's a dir, return the list of files
						files = fily.Content
					} else { // If it's a file, return it
						file = fily
						found = true
						break
					}
				}
			}
			if found {
				break
			}
		}

		// Throw an error to the client
		if !found {
			http.NotFound(ctx.Writer, ctx.Request)
			return
		}
	} else {
		// If it's a single file ("a.html") just check in the proj list
		for _, fily := range projecte.Files {
			if fily.Name == ctx.Params["filepath"] {
				file = fily
				found = true
				break
			}
		}

		// Throw an error if not found
		if !found {
			http.NotFound(ctx.Writer, ctx.Request)
			return
		}
	}

	// Read the file
	content, err := os.ReadFile(file.Path)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
	}

	// Gather the information relative to the file
	langName, loc := language.GetLangInfo(file)
	information := map[string]any{
		"lines": len(strings.Split(string(content), "\n")),
		"LOC":   loc,
		"lang":  langName,
	}

	cookies, err := makeCookieMap(ctx)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	// Unify data
	data := SingleFilePage{
		Name:     projecte.Name,
		Id:       projecte.Id,
		Files:    projecte.Files,
		Path:     file.Path,
		Content:  render.StringToChroma(string(content)),
		FileName: file.Name,
		Info:     information,
		Cookies:  cookies,
	}

	// Send the template
	ctx.Template([]string{"website/singlefile.html", "website/dirfileshower.html"}, data, nil)
}
