package handler

import (
	"diskhub/web/models"
	"diskhub/web/search"
	"net/http"
	"time"

	"github.com/esmyxvatu/feather"
)

type SearchPage struct {
	SearchTerm string
	Projects   []models.Project
	Time       int64
	Cookies    map[string]any
}

func SearchHandler(ctx *feather.Context) {
	// Get results and mesures the time taken for the search
	start := time.Now().UnixMilli()
	results, err := search.SearchFor(ctx.Query("q"))
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}
	end := time.Now().UnixMilli()

	cookies, err := makeCookieMap(ctx)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	// Unify all data in one variable
	data := SearchPage{
		SearchTerm: ctx.Query("q"),
		Projects:   results,
		Time:       end - start,
		Cookies:    cookies,
	}

	// Send the template
	ctx.Template([]string{"website/search.html"}, data, nil)
}
