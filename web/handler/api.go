package handler

import (
	"diskhub/web/filesystem"
	"diskhub/web/language"
	"diskhub/web/logger"
	"diskhub/web/models"
	"diskhub/web/ollama"
	"diskhub/web/wiki"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/esmyxvatu/feather"
)

func APISaveHandler(ctx *feather.Context) {
	id := ctx.FormValue("id")

	// Modify the project
	for i := range models.Projects {
		if models.Projects[i].Id == id {
			models.Projects[i].Saved = !models.Projects[i].Saved

			file, err := os.OpenFile(models.Projects[i].Path+"/.diskhub/.id", os.O_WRONLY, os.ModeAppend)
			logger.Console.Verify(err)
			defer file.Close()

			_, err = file.WriteString(id + "\ntrue")
			logger.Console.Verify(err)

			break
		}
	}

	// Redirect the user
	ctx.Redirect(http.StatusSeeOther, "/project/"+id)
}

func APIRefreshHandler(ctx *feather.Context) {
	// Reinitialize the global variables
	language.Langs = []language.Language{}
	models.Projects = []models.Project{}
	language.LangStats = []language.IndexLangs{}

	// Reload the different global variables
	language.InitLangs()
	filesystem.InitProjects()
	language.InitStats()

	// Redirect the user
	ctx.Redirect(http.StatusSeeOther, "/")
}

func APIReloadWiki(ctx *feather.Context) {
	id := ctx.Query("projID")

	for _, project := range models.Projects {
		if project.Id == id {
			articles, pathMap, err := wiki.Generate(project.Wiki.OriginDir)
			if err != nil {
				ctx.Error(http.StatusInternalServerError, err.Error())
			}

			project.Wiki = models.Wiki{Articles: articles, PathMap: pathMap, OriginDir: project.Wiki.OriginDir}
		}
	}

	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/doc/%s/%s", id, ctx.Query("origin")))
}

func APISetCookie(ctx *feather.Context) {
	//==================== Theme ====================================
	theme := ctx.Query("theme")

	if theme != "dark" && theme != "light" {
		ctx.Error(http.StatusBadRequest, "Invalid theme")
		return
	}

	// Define the cookie
	cookie := &http.Cookie{
		Name:    "theme",
		Value:   theme,
		Path:    "/",
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}

	// Set the cookie and redirect the user
	ctx.SetCookie(cookie)
	ctx.Redirect(http.StatusSeeOther, "/")
}

func APIAskOllama(ctx *feather.Context) {
	var data ollama.UserQuestion
	err := ctx.JSONBody(&data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ollama.UIMessage{Role: "system", Content: template.HTML(fmt.Sprintf("An error occured: %s", err.Error()))})
	}

	answer := ollama.AskOllama(data.Model, data.Content, data.Stream)

	ctx.JSON(http.StatusOK, answer)
}
