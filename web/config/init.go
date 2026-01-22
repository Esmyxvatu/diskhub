package config

import (
	"diskhub/web/logger"
	"diskhub/web/models"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func InitConfig() MainConfig {
	var configPath string = "config.toml"
	content, err := os.ReadFile(configPath)
	logger.Console.Verify(err)

	var cfg MainConfig
	err = toml.Unmarshal([]byte(content), &cfg)
	logger.Console.Verify(err)

	cfg.Diskhub.Exclude = append(cfg.Diskhub.Exclude, models.ExcludeList...)

	return cfg
}
