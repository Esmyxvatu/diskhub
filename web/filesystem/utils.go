package filesystem

import (
	"crypto/rand"
	"diskhub/web/models"
	"encoding/hex"
)

func createId(length int16) (string, error) {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
func addId(listdir []models.FileObject, id string, path string) []models.FileObject {
	for i := range listdir {
		listdir[i].ProjID = id

		// Use recursively the method if it's a dir
		if listdir[i].IsDir {
			listdir[i].Content = addId(listdir[i].Content, id, path+"/"+listdir[i].Name)
		} else {
			listdir[i].ProjPath = path + "/" + listdir[i].Name
		}
	}

	return listdir
}
