package filesystem

import (
	"diskhub/web/config"
	"diskhub/web/logger"
	"diskhub/web/models"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
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

func CreateExcludeList(excludeFilePath string) ([]string, error) {
	excludeRegex := []string{}

	file, err := os.Open(excludeFilePath)
	if err != nil {
		return nil, err
	}

	byteContent := make([]byte, 1000*1024)
	_, err = file.Read(byteContent)
	if err != nil {
		return nil, err
	}

	content := string(byteContent)

	lines := strings.SplitSeq(content, "\n")
	for rawLine := range lines {
		line := strings.TrimSpace(rawLine)

		if len(line) > 50*1024 || line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "use_gitignore") && strings.Split(line, "=")[1] == "true" {
			toExclude, err := CreateExcludeList(".gitignore")
			if err != nil {
				return nil, err
			}

			excludeRegex = append(excludeRegex, toExclude...)
		}

		line = regexp.QuoteMeta(line)
		line = strings.ReplaceAll(line, "\\*", "([^/]+)")

		regex := "^" + line + "$"

		excludeRegex = append(excludeRegex, regex)
	}

	return excludeRegex, nil
}
