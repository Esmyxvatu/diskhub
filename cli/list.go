package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

type Config struct {
	Title  string   `json:"name"`
	Auhtor string   `json:"author"`
	About  string   `json:"desc"`
	Tags   []string `json:"tags"`
	Status string   `json:"status"`
}

type Boxes struct {
	Height int
	Rows   []string
}

func ListProjects(args []string) {
	projectWidth := 50.0
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	projects := make([]Boxes, 0)

	for _, arg := range args {
		files, err := os.ReadDir(arg)

		if err != nil {
			fmt.Printf("There's a problem while opening %s\n", arg)
			continue
		}

		for _, file := range files {
			if file.IsDir() {
				project, err := GetProject(file)
				if err != nil && err.Error() != "Empty Project" {
					fmt.Printf("A problem occured while parsing %s: %s \n", file.Name(), err.Error())
					continue
				} else if err == nil {
					// Make the box and then put it in projects
					rows := make([]string, 0)
					formattedText := fmt.Sprintf("\033[1;3m%s \033[90m%s\033[0m", project.Name, project.Path)
					rows = append(rows, fmt.Sprintf("%-65s", formattedText))

					nbSegment := int(math.Ceil(float64(len(project.Desc)) / (projectWidth - 2.0)))

					start := 0
					for i := range nbSegment {
						end := int(projectWidth) * (i + 1)

						if end > len(project.Desc) {
							end = len(project.Desc)
						} else if project.Desc[end - 1] != ' ' {
							for {
								if project.Desc[end - 1] == ' ' {
									break
								}
								end -= 1
							}
						}

						rows = append(rows, fmt.Sprintf("%-50s", project.Desc[start:end]))
						start = end
					}

					lines := "\033[96m"
					offset := 0
					for _, tag := range project.Tags {
						offset += len(tag) + 1

						if offset > int(projectWidth-2.0) {
							lines += strings.Repeat(" ", int(projectWidth) - offset)
							offset = 0

							lines += "\033[0m"
							rows = append(rows, lines)

							lines = "\033[96m"
						}

						lines += fmt.Sprint(tag + " ")
					}
					lines += strings.Repeat(" ", int(projectWidth) - offset) + "\033[0m"
					rows = append(rows, lines)
					rows = append(rows, "")

					projects = append(projects, Boxes{
						Rows:   rows,
						Height: len(rows),
					})
				}
			}
		}
	}

	nbColumn := math.Floor(float64(termWidth) / (projectWidth + 1))

	if nbColumn < 1 {
		nbColumn = 1
	}

	col := int(nbColumn)

	for i := 0; i < len(projects); i += col {
		end := i + col
		if end > len(projects) {
			end = len(projects)
		}

		lineProjects := projects[i:end]

		// Trouve la hauteur max parmi les boxes de cette ligne
		maxH := 0
		for _, p := range lineProjects {
			if p.Height > maxH {
				maxH = p.Height
			}
		}

		// Affiche ligne par ligne
		for row := 0; row < maxH; row++ {
			for _, p := range lineProjects {
				if row < len(p.Rows) {
					fmt.Print(p.Rows[row])
				} else {
					// Si la box est plus petite, on remplit avec des espaces
					fmt.Print(strings.Repeat(" ", int(projectWidth)))
				}

				fmt.Print(" ") // petit espace entre chaque box
			}
			fmt.Println()
		}
	}

}

func GetProject(dir os.DirEntry) (Project, error) {
	var project Project

	absolutePath, err := filepath.Abs(dir.Name())
	if err != nil {
		return project, errors.New("Unable to get back the full path of folder")
	}

	files, err := os.ReadDir(absolutePath)
	if err != nil {
		return project, errors.New("Unable to open directory " + absolutePath)
	}

	for _, file := range files {
		if file.Name() == ".diskhub" && file.IsDir() {
			content, _ := os.ReadFile(absolutePath + "/.diskhub/conf.json")
			var cfg Config
			json.Unmarshal(content, &cfg)

			project.Path = absolutePath
			project.Name = cfg.Title
			project.Desc = cfg.About
			project.Tags = cfg.Tags

			return project, nil
		}
	}

	return project, errors.New("Empty Project")
}
