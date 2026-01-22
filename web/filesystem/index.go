package filesystem

import (
	"diskhub/web/config"
	"diskhub/web/logger"
	"diskhub/web/models"
	"os"
	"path/filepath"
	"slices"
)

func IndexDirectory(dirPath string) []models.FileObject {
	list := []models.FileObject{}

	// Get every files in the directory requested
	files, err := os.ReadDir(dirPath)
	logger.Console.Verify(err)

	for _, file := range files {
		// Ignore the automatic build dirs
		if slices.Contains(config.Configuration.Diskhub.Exclude, file.Name()) {
			continue
		}

		// Gather the information of the file
		fullPath := filepath.Join(dirPath, file.Name())
		infoFile, err := file.Info()
		logger.Console.Verify(err)
		modTime := infoFile.ModTime().Local().Format("2006-01-02 15:04:05")

		// Add the file to the list
		if file.IsDir() {
			actual := models.FileObject{Name: file.Name(), Path: fullPath, IsDir: true, Content: IndexDirectory(fullPath), LastMod: modTime}
			list = append(list, actual)
		} else {
			list = append(list, models.FileObject{Name: file.Name(), Path: fullPath, IsDir: false, Content: []models.FileObject{}, LastMod: modTime})
		}
	}

	return list
}
