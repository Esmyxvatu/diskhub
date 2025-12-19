package main

import (
	"regexp"
	"slices"
	"sort"
	"strings"
)

//============================================================ Type Definitions =====================================================================

type KeyValue struct {
	Project *Project
	Value   int
}

//============================================================ Functions ============================================================================

func SearchFor(query string) ([]Project, error) {
	var projectFound map[*Project]int = make(map[*Project]int)
	var params map[string]string = map[string]string{}

	// Create the regex with every words in the query
	keywords := strings.Split(query, " ")

	for _, kw := range keywords {
		// Check if it's a specifier (lang, tag, status)
		if strings.Contains(kw, ":") {
			temp := strings.Split(kw, ":")
			params[temp[0]] = temp[1]

			// Remove the keyword from the keywords list
			keywords = RemoveString(keywords, kw)
		}
	}

	// Create the regex for every keyword
	var regexPattern string
	for _, keyword := range keywords {
		if regexPattern != "" {
			regexPattern += "|"
		}
		regexPattern += regexp.QuoteMeta(strings.ToLower(keyword))
	}
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return []Project{}, err
	}

	// Search trough the whole list of project
	for _, proj := range projects {
		words := strings.Split(strings.ToLower(proj.About), " ")

		// Search for every occurence of the query (splited in word) in the Name
		matches := re.FindAllString(strings.ToLower(proj.Name), -1)
		for range matches {
			projectFound[&proj] += 1
		}

		// Search for every occurence of the query (splited in word) in the About
		for _, word := range words {
			matches := re.FindAllString(word, -1)
			for range matches {
				projectFound[&proj] += 1
			}
		}
	}

	// Organize Project in the order of points they have
	var kvPairs []KeyValue
	for project, value := range projectFound {
		if value > 0 {
			kvPairs = append(kvPairs, KeyValue{Project: project, Value: value})
		}
	}
	sort.Slice(kvPairs, func(i, j int) bool {
		return kvPairs[i].Value > kvPairs[j].Value
	})
	var OrganizedProjects []Project
	for _, kv := range kvPairs {
		OrganizedProjects = append(OrganizedProjects, *kv.Project)
	}

	// Check and filter with the specifier
	for key, value := range params {
		badIds := []string{}
		reverse := false
		if (value[0] == '!') { reverse = true; value = value[1:] }

		if key == "lang" {
			// Check if the project have the specified language
			for i := range len(OrganizedProjects) {
				proj := OrganizedProjects[i]
				langs := make([]string, len(proj.Langs))
				for i, lang := range proj.Langs {
					langs[i] = lang.Key
				}
				if !slices.Contains(langs, value) && !reverse {
					badIds = append(badIds, proj.Id)
				} else if slices.Contains(langs, value) && reverse {
					badIds = append(badIds, proj.Id)
				}
			}
		} else if key == "tag" {
			// Search if the project has the tag wanted
			for _, proj := range OrganizedProjects {
				if !slices.Contains(proj.Tags, value) && !reverse {
					badIds = append(badIds, proj.Id)
				} else if slices.Contains(proj.Tags, value) && reverse {
					badIds = append(badIds, proj.Id)
				}
			}
		} else if key == "status" {
			// Check if the project have the wanted status
			for _, proj := range OrganizedProjects {
				if !strings.Contains(strings.ToLower(proj.Status), strings.ToLower(value)) && !reverse {
					badIds = append(badIds, proj.Id)
				} else if strings.Contains(strings.ToLower(proj.Status), strings.ToLower(value)) && reverse {
					badIds = append(badIds, proj.Id)
				}
			}
		}

		// Remove the unwanted projects
		for _, id := range badIds {
			for i := range OrganizedProjects {
				if OrganizedProjects[i].Id == id {
					OrganizedProjects = slices.Delete(OrganizedProjects, i, i+1)
					break
				}
			}
		}
	}

	return OrganizedProjects, nil
}

func RemoveString(list []string, s string) []string {
	// Get the pos of the string
	for i, obj := range list {
		if obj == s {
			// Remove the string
			return slices.Delete(list, i, i+1)
		}
	}
	return list
}
