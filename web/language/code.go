package language

import (
	"bufio"
	"diskhub/web/models"
	"diskhub/web/logger"
	"os"
	"path/filepath"
	"strings"
)

func GetLangInfo(file models.FileObject) (string, float32) {
	// Found the type using extension
	ext := strings.TrimPrefix(filepath.Ext(file.Name), ".")

	lang := LangByExt[ext]
	if lang == nil {
		lang = &Language{Name: "Other"}
	}

	var nbLoc float32

	// Open the file and read it line by line
	fileObj, err := os.Open(file.Path)
	if err != nil {
		logger.Console.Fatal("%s", err.Error())
	}
	defer fileObj.Close()
	scanner := bufio.NewScanner(fileObj)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	// Get line after line and count them
	inBlock := false
	blockEnd := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if inBlock {
			if strings.Contains(line, blockEnd) {
				inBlock = false
			}
			continue
		}

		isComment := false

		for _, com := range lang.Coms {
			if com.End == "" {
				// commentaire ligne
				if strings.HasPrefix(line, com.Start) {
					isComment = true
					break
				}
			} else {
				// commentaire bloc
				if strings.HasPrefix(line, com.Start) || strings.Contains(line, com.Start) {
					inBlock = true
					blockEnd = com.End
					isComment = true
					break
				}
			}
		}

		if !isComment {
			nbLoc++
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Console.Error("%s", err.Error())
	}

	return lang.Name, nbLoc
}

func GetLangsStats(files []models.FileObject) []LanguageStat {
	var langs []LanguageStat

	// Loop over every file
	for _, file := range files {
		// If it's a dir, get the stats of it's content
		if file.IsDir {
			langy := GetLangsStats(file.Content)

			// Add the stats of the content to the stats of the root
			for i := range langy {
				name := langy[i].Name
				loc := langy[i].Loc
				file := langy[i].File
				found := false

				for j := range langs {
					if langs[j].Name == name {
						langs[j].Loc += loc
						langs[j].File += file
						found = true
						break
					}
				}

				// Create the data if it doesn't exist
				if !found {
					langs = append(langs, LanguageStat{
						Name: name,
						Loc:  loc,
						File: file,
					})
				}
			}
		} else {
			// Get the info of the file
			name, loc := GetLangInfo(file)
			found := false

			// Add the info to the data
			for i := range langs {
				if langs[i].Name == name {
					langs[i].Loc += loc
					langs[i].File++
					found = true
					break
				}
			}

			// Create the data if it doesn't exist
			if !found {
				langs = append(langs, LanguageStat{
					Name: name,
					Loc:  loc,
					File: 1,
				})
			}
		}
	}

	return langs
}

func RecursiveGetLang(dir models.FileObject) (map[string]float32, int) {
	var LangsInfo map[string]float32 = make(map[string]float32)
	var LinesOfCode int

	for _, file := range dir.Content {
		// If it's a file, add the info to the root stats
		// Else, do it recursively
		if !file.IsDir {
			langName, nbLoc := GetLangInfo(file)
			LangsInfo[langName] += nbLoc
			LinesOfCode += int(nbLoc)
		} else {
			langsInfo, nbLoc := RecursiveGetLang(file)
			LinesOfCode += nbLoc

			for key, value := range langsInfo {
				LangsInfo[key] += value
			}
		}
	}

	return LangsInfo, LinesOfCode
}
