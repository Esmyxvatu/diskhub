package filesystem

import (
	"diskhub/web/language"
	"diskhub/web/logger"
	"diskhub/web/models"
	"diskhub/web/wiki"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func EstablishProject(dir models.FileObject) (models.Project, error) {
	// Search for the .diskhub folder who contains all informations
	mainDir, err := dir.Found(".diskhub")
	if err != nil {
		return models.Project{}, err
	}
	// Open the file who declare all important information
	infoFile, err := mainDir.Found("conf.json")
	if err != nil {
		return models.Project{}, err
	}

	// Load the information in the project
	var Projecte models.Project
	file, err := os.Open(infoFile.Path)
	logger.Console.Verify(err)
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&Projecte); err != nil {
		return models.Project{}, err
	}

	// Create a unique ID for every Project and save it in .id (Load if .id exist)
	var id string
	saved := false
	idFile, err := mainDir.Found(".id")
	if err != nil && err != models.ErrFileNotFound {
		// Return the error if one happens
		return models.Project{}, err
	} else if err != nil {
		// Create a custom ID if the file isn't found
		id, err = createId(16)
		logger.Console.Verify(err)

		// Create the .id file to store the id
		file, err := os.Create(filepath.Join(mainDir.Path, ".id"))
		logger.Console.Verify(err)
		defer file.Close()

		// Append the id to the .id file (so it stays the same)
		_, err = file.WriteString(id)
		logger.Console.Verify(err)
	} else {
		// Open the .id file
		file, err := os.Open(idFile.Path)
		logger.Console.Verify(err)

		// Read the 16 first byte (= the length of the id)
		buffer := make([]byte, 16)
		_, err = file.Read(buffer)
		logger.Console.Verify(err)

		// Convcert the raw bytes to a usable string
		id = string(buffer)

		// Get the next 5 bytes (= wether the project is marked as "saved" or not)
		buffer = make([]byte, 5)
		n, err := file.Read(buffer)

		if err != nil && err.Error() != "EOF" {
			// An error occured while reading that ain't reaching the end of the file
			logger.Console.Error("Error while reading saved status: %v", err)
		}

		// Verify if the status is saved and mark the project accordingly
		if n == 5 && string(buffer) == "\ntrue" {
			saved = true
		}
	}

	// Get the .ignore file and create all corresponding regexes
	excludeFile, err := mainDir.Found(".ignore")
	excludeRegex := []string{"\\.diskhub"}
	if err != nil && err != models.ErrFileNotFound {
		// Return the error if one occurs
		return models.Project{}, err
	} else if err != nil {
		// If the file isn't found, create it and initialize it with the bare minimum
		file, err := os.Create(filepath.Join(mainDir.Path, ".ignore"))
		logger.Console.Verify(err)
		defer file.Close()

		_, err = file.WriteString("use_gitignore = false")
		logger.Console.Verify(err)
	} else {
		toExclude, err := CreateExcludeList(excludeFile.Path)
		logger.Console.Verify(err)

		excludeRegex = append(excludeRegex, toExclude...)
	}

	// Get the content of the models.Project, the langs used, and the number of lines of code
	var ProjContent []models.FileObject
	var LangsInfo map[string]float32 = make(map[string]float32)
	var LinesOfCode int
	for _, file := range dir.Content {
		toAvoid := false

		for _, regex := range excludeRegex {
			re := regexp.MustCompile(regex)
			logger.Console.Info("%v : %s", re, strings.TrimPrefix(file.Path, dir.Path+"/"))
			if re.MatchString(strings.TrimPrefix(file.Path, dir.Path+"/")) {
				toAvoid = true
				break
			}
		}

		if toAvoid {
			continue
		}

		ProjContent = append(ProjContent, file)

		// Load the different information, recursively for the dirs
		if !file.IsDir {
			langName, nbLoc := language.GetLangInfo(file)
			LangsInfo[langName] += nbLoc
			LinesOfCode += int(nbLoc)
		} else {
			langsInfo, nbLoc := language.RecursiveGetLang(file)
			LinesOfCode += nbLoc

			for key, value := range langsInfo {
				LangsInfo[key] += value
			}
		}
	}

	// Calculate the percentage of use of every languages
	for key, value := range LangsInfo {
		LangsInfo[key] = value / float32(LinesOfCode) * 100
	}

	// Sort the the percentage by top to down
	var sortedLangs []models.StringInt
	for k, v := range LangsInfo {
		sortedLangs = append(sortedLangs, models.StringInt{Key: k, Value: v, IsLast: false})
	}
	sort.Slice(sortedLangs, func(i, j int) bool {
		return sortedLangs[i].Value > sortedLangs[j].Value
	})

	// Load the ReadMe file
	readmeFile, err := dir.Found("README.md")
	if err != nil && err != models.ErrFileNotFound {
		logger.Console.Error("%s", err.Error())
	}

	wiki, pathMap, err := wiki.Generate(mainDir)
	if err != nil && err != models.ErrFileNotFound {
		logger.Console.Error("%s", err.Error())
	}

	// Define the models.Project and return it
	Projecte.Id = id
	Projecte.Files = addId(ProjContent, id, "")
	Projecte.Langs = sortedLangs
	Projecte.LOC = LinesOfCode
	Projecte.ReadMe = readmeFile.Path
	Projecte.LastMod = dir.GetLastUpdate()
	Projecte.Saved = saved
	Projecte.Wiki = models.Wiki{Articles: wiki, PathMap: pathMap, OriginDir: dir}
	Projecte.Path = dir.Path

	return Projecte, nil
}
