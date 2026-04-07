package main

import (
	"fmt"

	"diskhub/web/config"
	"diskhub/web/filesystem"
	"diskhub/web/handler"
	"diskhub/web/language"
	"diskhub/web/logger"
	"diskhub/web/snippet"

	"github.com/esmyxvatu/feather"
	"github.com/esmyxvatu/feather/middlewares"
)

func main() {
	server := feather.NewServer()
	server.AddMiddleware(
		middlewares.Logging(),
	)

	// Deserve static files
	server.Static("/public", "website/public")

	// User Routes
	server.GET("/", handler.IndexHandler)
	server.GET("/settings", handler.SettingsHandler)
	server.GET("/search", handler.SearchHandler)
	server.GET("/project/:id|[a-f0-9]+", handler.ProjectHandler)
	server.GET("/project/:id|[a-f0-9]+/tree/*filepath", handler.FileShowerHandler)
	server.GET("/doc/:id|[a-f0-9]+/*page", handler.DocShowerHandler)
	server.GET("/file-explorer/*path", handler.FileExplorerHandler)
	server.GET("/snippet", handler.SnippetHandler)
	if config.Configuration.Ollama.Active {
		server.GET("/chat", handler.ChatHandler)
	}

	// API Routes
	server.POST("/api/save", handler.APISaveHandler)
	server.POST("/api/refresh", handler.APIRefreshHandler)
	server.POST("/api/cookie", handler.APISetCookie)
	server.POST("/api/askollama", handler.APIAskOllama)
	server.POST("/api/reloadwiki", handler.APIReloadWiki)
	server.POST("/api/toggleOllama", handler.APIToggleOllama)

	// Admin Routes
	server.GET("/admin/dashboard", handler.AdminDashboardHandler)

	// Setup the projects, langs and stats
	language.InitLangs()
	filesystem.InitProjects()
	language.InitStats()
	snippet.InitDb()

	// Run the server
	logger.Console.Info("Starting server on port %d", config.Configuration.Server.Port)
	server.Listen(fmt.Sprintf(":%d", config.Configuration.Server.Port))
}
