package models

import (
	"diskhub/web/logger"
	"time"
)

type FileObject struct {
	Name     string
	Path     string
	IsDir    bool
	LastMod  string
	Content  []FileObject
	ProjID   string
	ProjPath string
}

func (dir FileObject) GetLastUpdate() string {
	// Define the minimum to the dir last mod
	result, err := time.Parse("2006-01-02 15:04:05", dir.LastMod)
	if err != nil {
		logger.Console.Fatal("%s", err.Error())
	}

	for _, file := range dir.Content {
		// If it's a dir, apply recursively the function
		if file.IsDir {
			date, err := time.Parse("2006-01-02 15:04:05", file.GetLastUpdate())
			logger.Console.Verify(err)

			if date.After(result) {
				result = date
			}
		} else {
			date, err := time.Parse("2006-01-02 15:04:05", file.LastMod)
			logger.Console.Verify(err)

			if date.After(result) {
				result = date
			}
		}
	}

	return result.Format("2006-01-02 15:04:05")
}

func (dir FileObject) Found(name string) (FileObject, error) {
	// Search for a specific file in a dir
	for _, item := range dir.Content {
		if item.Name == name {
			return item, nil
		}
	}

	return FileObject{}, ErrFileNotFound
}
