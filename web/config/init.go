package config

import (
	"diskhub/web/logger"
	"diskhub/web/models"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func InitConfig() (MainConfig, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return MainConfig{}, err
	}

	var cfg MainConfig
	err = toml.Unmarshal([]byte(content), &cfg)
	if err != nil {
		return MainConfig{}, err
	}

	cfg.Diskhub.Exclude = append(cfg.Diskhub.Exclude, models.ExcludeList...)

	return cfg, nil
}

func MustInitConfig() MainConfig {
	cfg, err := InitConfig()
	if err != nil {
		logger.Console.Fatal("An error occured while parsing config file: %s", err.Error())
	}

	return cfg
}