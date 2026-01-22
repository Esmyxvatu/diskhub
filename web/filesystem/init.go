package filesystem

import (
	"diskhub/web/config"
	"diskhub/web/logger"
	"diskhub/web/models"
)

func InitProjects() {
	// Loop on the models.Projects global variables
	for _, path := range config.Configuration.Diskhub.Paths {
		dirs := IndexDirectory(path)
		for _, item := range dirs {
			// If it found the diskhub.toml file, add it, else ignore the dir
			Projecte, err := EstablishProject(item)
			if err != nil {
				if err == models.ErrFileNotFound {
					continue
				}
			} else {
				logger.Console.Info("Project %s found in %s", Projecte.Name, item.Path)
				models.Projects = append(models.Projects, Projecte)
			}
		}
	}
}
