package handler

import (
	"diskhub/web/models"

	"github.com/esmyxvatu/feather"
)

func makeCookieMap(ctx *feather.Context) (map[string]any, error) {
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err != nil {
		return map[string]any{}, models.ErrParsingCookies
	}

	if themeCookie.Value == "dark" {
		dark = true
	}

	// Create the cookies dict
	cookies := map[string]any{
		"Theme": dark,
	}

	return cookies, nil
}
