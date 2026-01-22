package config

/*
	Code representation of the config.toml file.
	Englob the config of Diskhub (path to search for projects and file to exclude) and Ollama (wether it's active or not).
*/
var Configuration MainConfig = InitConfig()