package main

import (
	"time"
	"os"
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"github.com/esmyxvatu/feather"
	"github.com/esmyxvatu/feather/middlewares"
)

//============================================================ Struct Types =========================================================================
type ServerConfig struct {
	Address string
	Port 	int
}
type DiskhubConfig struct {
	Paths  []string
	Exclude []string
}
type OllamaConfig struct {
	Active bool
	Port   int
}
type MainConfig struct {
	Server  ServerConfig
	Diskhub DiskhubConfig
	Ollama  OllamaConfig
}

//============================================================ Global Variables =====================================================================

var projects 	 []Project
var langStat 	 []IndexLangs
var console      *Logger          = &Logger{}
var Uptime_start time.Time 	  = time.Now()
var Configuration MainConfig      = InitConfig()

//============================================================ Functions ============================================================================

func InitConfig() MainConfig {
	var configPath string = "config.toml"
	content, err := os.ReadFile(configPath)
	console.verify(err)

	var cfg MainConfig
	err = toml.Unmarshal([]byte(content), &cfg)
	console.verify(err)
	
	cfg.Diskhub.Exclude = append(cfg.Diskhub.Exclude, ExcludeList...)

	return cfg
}

func main() {
	server := feather.NewServer()
	server.AddMiddleware(
		middlewares.Logging(),
	)
	
	// Deserve static files
	server.Static("/public", "website/public")

	// User Routes
	server.GET("/", IndexHandler)
	server.GET("/search", SearchHandler)
	server.GET("/project/:id|[a-f0-9]+", ProjectHandler)
	server.GET("/project/:id|[a-f0-9]+/tree/*filepath", FileShowerHandler)
	server.GET("/project/:id|[a-f0-9]+/doc/:page|[a-zA-Z0-9._-]+", DocShowerHandler)
	server.GET("/file-explorer/*path", FileExplorerHandler)
	server.GET("/snippet", SnippetHandler)
	if !Configuration.Ollama.Active {
		server.GET("/chat", ChatHandler)
	}

	// API Routes
	server.POST("/api/save", APISaveHandler)
	server.POST("/api/refresh", APIRefreshHandler)
	server.POST("/api/cookie", APISetCookie)
	server.POST("/api/askollama", APIAskOllama)

	// Admin Routes
	server.GET("/admin/dashboard", AdminDashboardHandler)

	// Setup the projects, langs and stats
	InitLangs()
	InitProjects()
	InitStats()
	InitSnippetDb()

	// Run the server
	console.info("Starting server on port %d", Configuration.Server.Port)
	server.Listen(fmt.Sprintf(":%d", Configuration.Server.Port))
}
