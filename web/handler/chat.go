package handler

import (
	"diskhub/web/ollama"

	"github.com/esmyxvatu/feather"
)

type ChatPage struct {
	Models   []string
	Messages []ollama.UIMessage
	Cookies  map[string]any
}

func ChatHandler(ctx *feather.Context) {
	// Check for the theme cookie
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

	// Unify the data
	data := ChatPage{
		Models:   ollama.ListModelAvaillable(),
		Messages: ollama.UIMessages,
		Cookies:  cookies,
	}

	// Send the template
	ctx.Template([]string{"website/chat.html"}, data, nil)
}
