package config


type ServerConfig struct {
	Address string `toml:"adress"`			// Address on wich the server should listen
	Port    int    `toml:"port"`			// Port the server should use
}

type DiskhubConfig struct {
	Paths   []string `toml:"paths"`			// Paths to search for projects
	Exclude []string `toml:"exclude"`		// Name of files/folders to automaticly exclude
}

type OllamaConfig struct {
	Active bool `toml:"activate"`			// If the user want ollama to be activated or not
	Port   int  `toml:"port"`				// On wich port ask for ollama response
}

type MainConfig struct {
	Server  ServerConfig  `toml:"Server"`	// Params related to the server
	Diskhub DiskhubConfig `toml:"Diskhub"`	// Params related to Diskhub
	Ollama  OllamaConfig  `toml:"Ollama"`	// Params related to Ollama
}