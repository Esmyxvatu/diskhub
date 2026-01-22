package language

import (
	"diskhub/web/logger"
	"diskhub/web/models"
	"encoding/json"
	"os"
	"sort"
)

func InitLangs() {
	// Read the languages config files
	content, err := os.ReadFile("languages.json")
	if err != nil {
		logger.Console.Fatal("%s", err.Error())
	}

	// Load the configuration
	json.Unmarshal(content, &Langs)

	// Define the maps of all Langs
	LangByExt = make(map[string]*Language)
	for i := range Langs {
		lang := &Langs[i]
		for _, ext := range lang.Exts {
			LangByExt[ext] = lang
		}
	}
}

func InitStats() {
	for _, proj := range models.Projects {
		langi := GetLangsStats(proj.Files) // Get the stats of every files

		// Unify them in one big Stats
		for _, lang := range langi {
			name := lang.Name
			found := false

			// If it's already found, add it to existing data
			for i := range LangStats {
				if LangStats[i].Name == name {
					LangStats[i].Loc += int(lang.Loc)
					LangStats[i].Files += lang.File
					LangStats[i].Projs++
					found = true
					break
				}
			}

			// Else, create the data
			if !found {
				LangStats = append(LangStats, IndexLangs{
					Name:  name,
					Loc:   int(lang.Loc),
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
	for _, stat := range LangStats {
		if IsCode, ok := filterMap[stat.Name]; ok && IsCode {
			LangStats[i] = stat
			i++
		}
	}
	LangStats = LangStats[:i]

	// Sort the Langs stats top to down
	sort.Slice(LangStats, func(i, j int) bool {
		return LangStats[i].Loc > LangStats[j].Loc
	})
}
