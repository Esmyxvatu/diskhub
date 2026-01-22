package wiki

import (
	"diskhub/web/models"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func OrganizeArticles(articles *[]any, object models.FileObject) error {
	orgFile, err := object.Found(".org")
	if err != nil && err != models.ErrFileNotFound {
		return err
	}
	if orgFile.IsDir {
		return fmt.Errorf(".org need to be a file in %s", object.Name)
	}
	content, err := os.ReadFile(orgFile.Path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	organization := make(map[string]int)
	for _, line := range lines {
		parts := strings.Split(line, ":")

		name := strings.TrimSpace(parts[1])
		value, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return err
		}

		organization[name] = value
	}

	return nil
}
