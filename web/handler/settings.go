package handler

import (
	"fmt"
	"net/http"

	"diskhub/web/config"

	"github.com/esmyxvatu/feather"
)

type SettingsPage struct {
	Params  map[string]any
	Cookies map[string]any
}

func SettingsHandler(ctx *feather.Context) {
	cookies, err := makeCookieMap(ctx)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	params := map[string]any{
		"Ollama": config.Configuration.Ollama.Active,
	}

	fmt.Println(config.Configuration.Ollama.Active)

	// Create the dict
	data := SettingsPage{
		Params:  params,
		Cookies: cookies,
	}

	ctx.Template([]string{"website/settings.html"}, data, nil)
}
