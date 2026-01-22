package handler

import (
	"diskhub/web/models"
	"diskhub/web/render"
	"html/template"
	"runtime"
	"slices"
	"time"

	"github.com/esmyxvatu/feather"
)

var Ram_Values []models.StringInt
var MaxRamValue int = 50
var Uptime_start time.Time = time.Now()

type DashboardPage struct {
	Uptime  time.Duration
	RAM     float32
	Cookies map[string]any
	Chart   template.HTML
}

func AdminDashboardHandler(ctx *feather.Context) {
	// Load the theme cookie
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err == nil {
		if themeCookie.Value == "dark" {
			dark = true
		}
	}

	// Create the cookies dict
	cookies := map[string]any{
		"Theme": dark,
	}

	// Load and read the stats for the usage of memory
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	Ram_Values = append(Ram_Values, models.StringInt{Key: "", Value: float32(m.Alloc / 1024), IsLast: false})
	if len(Ram_Values) > MaxRamValue {
		Ram_Values = slices.Delete(Ram_Values, 0, 1)
	}

	// Unify all the data
	data := DashboardPage{
		Uptime:  time.Since(Uptime_start),
		RAM:     Ram_Values[len(Ram_Values)-1].Value,
		Cookies: cookies,
		Chart:   template.HTML(render.GenerateLineChart(Ram_Values, MaxRamValue, 480, 200)),
	}

	// Send the template
	ctx.Template([]string{"website/dashboard.html"}, data, nil)
}
