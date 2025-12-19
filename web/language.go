package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"sort"
	"slices"
)

//============================================================ Type Definitions =====================================================================

type Commentaire struct {
	Start string
	End   string
}
type Language struct {
	IsCode bool
	Name   string
	Exts   []string
	Coms   []Commentaire
}
type LangStat struct {
	Name string
	Loc  float32
	File int
}

//============================================================ Global Variables =====================================================================

var Langs []Language

//============================================================ Functions ============================================================================

func InitLangs() {
	// Read the languages config files
	content, err := os.ReadFile("languages.json")
	if err != nil {
		console.fatal("%s", err.Error())
	}

	// Load the configuration
	json.Unmarshal(content, &Langs)
}
func InitStats() {
	for _, proj := range projects {
		langi := GetLangsStats(proj.Files)	// Get the stats of every files

		// Unify them in one big Stats
		for _, lang := range langi {
			name := lang.Name
			found := false

			// If it's already found, add it to existing data
			for i := range langStat {
				if langStat[i].Name == name {
					langStat[i].Loc += int(lang.Loc)
					langStat[i].Files += lang.File
					langStat[i].Projs++
					found = true
					break
				}
			}

			// Else, create the data
			if !found {
				langStat = append(langStat, IndexLangs{
					Name: name,
					Loc: int(lang.Loc),
					Projs: 1,
					Files: lang.File,
				})
			}
		}
	}

	filterMap := make(map[string]bool)
	for _, lang := range Langs {
		filterMap[lang.Name] = lang.IsCode
	}

	i := 0
	for _, stat := range langStat {
		if IsCode, ok := filterMap[stat.Name]; ok && IsCode {
			langStat[i] = stat
			i++
		}
	}
	langStat = langStat[:i]

	// Sort the Langs stats top to down
	sort.Slice(langStat, func(i, j int) bool {
		return langStat[i].Loc > langStat[j].Loc
	})
}
func (file FileObject) GetLangInfo() (string, float32) {
	// Found the type using extension
	exts := strings.Split(file.Name, ".")
	ext := exts[len(exts)-1]

	var lang Language = Language{ Name: "Other" }
	for _, languy := range Langs {
		if slices.Contains(languy.Exts, ext) {
			lang = languy
			break
		}
	}

	var nbLoc float32

	// Open the file and read it line by line
	fileObj, err := os.Open(file.Path)
	if err != nil {
		console.fatal("%s", err.Error())
	}
	defer fileObj.Close()
	scanner := bufio.NewScanner(fileObj)

	// Get line after line and count them
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || line == "\n" {
			continue
		}

		commented := false
		for _, com := range lang.Coms {
			if strings.HasPrefix(line, com.Start) {
				commented = true
			} else if strings.Contains(line, com.Start) {
				commented = true
				nbLoc++
			}

			if strings.HasSuffix(line, com.End) {
				commented = false
			} else if strings.Contains(line, com.End) {
				commented = false
				nbLoc++
			}
		}

		if !commented { nbLoc++ }
	}
	if err := scanner.Err(); err != nil {
		console.error("%s", err.Error())
	}

	return lang.Name, nbLoc
}
func GetLangsStats(files []FileObject) []LangStat {
	var langs []LangStat

	// Loop over every file
	for _, file := range files {
		// If it's a dir, get the stats of it's content
		if file.IsDir {
			langy := GetLangsStats(file.Content)

			// Add the stats of the content to the stats of the root
			for i := range langy {
				name := langy[i].Name
				loc  := langy[i].Loc
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
					langs = append(langs, LangStat{
						Name: name,
						Loc: loc,
						File: 1,
					})
				}
			}
		} else {
			// Get the info of the file
			name, loc := file.GetLangInfo()
			found := false

			// Add the info to the data
			for i := range langs {
				if langs[i].Name == name {
					langs[i].Loc += loc
					langs[i].File ++
					found = true
					break
				}
			}

			// Create the data if it doesn't exist
			if !found {
				langs = append(langs, LangStat{
					Name: name,
					Loc: loc,
					File: 1,
				})
			}
		}
	}

	return langs
}
func RecursiveGetLang(dir FileObject) (map[string]float32, int) {
	var LangsInfo map[string]float32 = make(map[string]float32)
	var LinesOfCode int

	for _, file := range dir.Content {
		// If it's a file, add the info to the root stats
		// Else, do it recursively
		if !file.IsDir {
			langName, nbLoc := file.GetLangInfo()
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
