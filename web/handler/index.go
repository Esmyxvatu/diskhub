package handler

import (
	"html/template"
	"net/http"
	"slices"

	"diskhub/web/language"
	"diskhub/web/models"
	"diskhub/web/render"

	"github.com/esmyxvatu/feather"
)

type IndexPage struct {
	Projects    []models.Project
	Info        map[string]any
	Langs       []language.IndexLangs
	StatusChart template.HTML
	Status      map[string]float64
	Cookies     map[string]any
}

func IndexHandler(ctx *feather.Context) {
	// Calculate the different information
	var nbfiles int
	var nbloc int
	for i := range language.LangStats {
		nbfiles += language.LangStats[i].Files
		nbloc += language.LangStats[i].Loc
	}
	infos := map[string]any{
		"nbLangs": len(language.LangStats),
		"nbFiles": nbfiles,
		"nbLOC":   nbloc,
	}

	// Get the status of every projects, and save it in a dict
	statusList := map[string]float64{}
	found := []string{}
	for _, proj := range models.Projects {
		if slices.Contains(found, proj.Status) {
			statusList[proj.Status]++
		} else {
			statusList[proj.Status] = 1
			found = append(found, proj.Status)
		}
	}

	cookies, err := makeCookieMap(ctx)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	// Create the dict
	data := IndexPage{
		Projects:    models.Projects,
		Info:        infos,
		Langs:       language.LangStats,
		StatusChart: template.HTML(render.GeneratePieChart(statusList, 300, 300)), // Escape the generated pie with template.HTML to be interpreted as HTML
		Status:      statusList,
		Cookies:     cookies,
	}

	// Send the template
	ctx.Template([]string{"website/index.html"}, data, nil)
}
