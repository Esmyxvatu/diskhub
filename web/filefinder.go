package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
)

//============================================================ Variables ============================================================================
var ExcludeList []string = []string{
	// Git
	".git",
	".gitignore",
	".gitattributes",

	// JavaScript - TypeScript
	"node_modules",
	"package-lock.json",
	"tsconfig.json",
	
	// Rust
	"target",
	"Cargo.lock",
	
	// Golang
	"go.mod",
	"go.sum",
	
	// Zig
	"zig-build",
	
	// Python
	"__pycache__",
	"requirements.txt",
	
	// Compiled programs
	"*.out",
	"*.exe",

	// Cpp
	".cache",
}

//============================================================ Type declarations ====================================================================

type FileObject struct {
	Name     string
	Path     string
	IsDir    bool
	LastMod  string
	Content  []FileObject
	ProjID   string
	ProjPath string
}
type GeneralConfig struct {
	Title  string
	Auhtor string
	About  string
	Tags   []string
	Status string
}
type ExceptConfig struct {
	Files []string
	Dirs  []string
}
type ProjectFileConfig struct {
	General GeneralConfig
	Except  ExceptConfig
	Links   map[string]string
	Doc     map[string]string
}
type Project struct {
	Name    string				`json:"name"`
	Id      string
	About   string				`json:"desc"`
	Links   map[string]string	`json:"links"`
	Files   []FileObject
	Langs   []StringInt
	LOC     int
	LastMod string
	doc     string
	Tags    []string			`json:"tags"`
	Doc     map[string]string
	Saved   bool
	Path    string
	Status  string				`json:"status"`
}
type StringInt struct {
	Key    string
	Value  float32
	IsLast bool
}

//============================================================ Functions ============================================================================

// 	Helper method
func CreateId(length int16) (string, error) {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
func addId(listdir []FileObject, id string, path string) []FileObject {
	for i := range listdir {
		listdir[i].ProjID = id

		// Use recursively the method if it's a dir
		if listdir[i].IsDir {
			listdir[i].Content = addId(listdir[i].Content, id, path + "/" + listdir[i].Name)
		} else {
			listdir[i].ProjPath = path + "/" + listdir[i].Name
		}
	}

	return listdir
}

//	Method for the FileObject Type
func (dir FileObject) GetLastUpdate() string {
	// Define the minimum to the dir last mod
	result, err := time.Parse("2006-01-02 15:04:05", dir.LastMod)
	if err != nil {
		console.fatal("%s", err.Error())
	}

	for _, file := range dir.Content {
		// If it's a dir, apply recursively the function
		if file.IsDir {
			date, err := time.Parse("2006-01-02 15:04:05", file.GetLastUpdate())
			if err != nil {
				console.fatal("%s", err.Error())
			}

			if date.After(result) {
				result = date
			}
		} else {
			date, err := time.Parse("2006-01-02 15:04:05", file.LastMod)
			if err != nil {
				console.fatal("%s", err.Error())
			}

			if date.After(result) {
				result = date
			}
		}
	}

	return result.Format("2006-01-02 15:04:05")
}
func (dir FileObject) EstablishProject() (Project, error) {
	// Search for the diskhub.toml file and read it's content
	mainDir, err := dir.Found(".diskhub")
	if err != nil {
		return Project{}, err
	}
	infoFile, err := mainDir.Found("conf.json")
	if err != nil {
		return Project{}, err
	}
	
	var projecte Project
	file, err := os.Open(infoFile.Path)
	console.verify(err)
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&projecte); err != nil {
		return Project{}, err
	}

	// Create a unique ID for every project and save it in .id (Load if .id exist)
	var id string
	saved := false
	idFile, err := mainDir.Found(".id")
	if err != nil && err.Error() != "FILE NOT FOUND" {
		return Project{}, err
	} else if err != nil {
		id, err = CreateId(16)
		console.verify(err)
		
		file, err := os.Create( filepath.Join(mainDir.Path, ".id") )
		console.verify(err)
		defer file.Close()

		_, err = file.WriteString(id)
		console.verify(err)
	} else {
		file, err := os.Open(idFile.Path)
		console.verify(err)
		
		buffer := make([]byte, 16)
		_, err = file.Read(buffer)
		console.verify(err)
		
		id = string(buffer)

		buffer = make([]byte, 5)
		n, err := file.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			console.error("Error while reading saved status: %v", err)
		}
		if n == 5 && string(buffer) == "\ntrue" {
			saved = true
		}
	}

	// Get the .ignore file and create all corresponding regexes
	excludeFile, err := mainDir.Found(".ignore")
	excludeRegex := []string{"\\.diskhub"}
	if err != nil && err.Error() != "FILE NOT FOUND" {
		return Project{}, err
	} else if err != nil {
		file, err := os.Create( filepath.Join(mainDir.Path, ".ignore") )
		console.verify(err)
		defer file.Close()

		_, err = file.WriteString("use_gitignore = false")
		console.verify(err)
	} else {
		file, err := os.Open(excludeFile.Path)
		console.verify(err)

		byteContent := make([]byte, 1000 * 1024)
		_, err = file.Read(byteContent)
		console.verify(err)
		content := string(byteContent)

		lines := strings.SplitSeq(content, "\n")
		for rawLine := range lines {
			line := strings.TrimSpace(rawLine)

			if len(line) > 50 * 1024 || line == "" || strings.HasPrefix(line, "#") { continue }
			if strings.Contains(line, "use_gitignore") {
				// TODO: Implement the support for Git ignore
				continue
			}

			line = regexp.QuoteMeta(line)
			line = strings.ReplaceAll(line, "\\*", "([^/]+)")

			regex := "^" + line + "$"

			excludeRegex = append(excludeRegex, regex)
		}
	}

	// Get the content of the project, the langs used, and the number of lines of code
	var ProjContent []FileObject
	var LangsInfo map[string]float32 = make(map[string]float32)
	var LinesOfCode int
	for _, file := range dir.Content {
		toAvoid := false

		for _, regex := range excludeRegex {
			re := regexp.MustCompile(regex)
			if re.MatchString( strings.TrimPrefix(file.Path, dir.Path + "/") ) {
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

	// Calculate the percentage of use of every languages
	for key, value := range LangsInfo {
		LangsInfo[key] = value / float32(LinesOfCode) * 100
	}

	// Sort the the percentage by top to down
	var sortedLangs []StringInt
	for k, v := range LangsInfo {
		sortedLangs = append(sortedLangs, StringInt{k, v, false})
	}
	sort.Slice(sortedLangs, func(i, j int) bool {
		return sortedLangs[i].Value > sortedLangs[j].Value
	})

	// Load the ReadMe file
	readmeFile, err := dir.Found("README.md")
	if err != nil && err.Error() != "FILE NOT FOUND" {
		console.error("%s", err.Error())
	}

	// Define the project and return it
	projecte.Id = id
	projecte.Files = addId(ProjContent, id, "")
	projecte.Langs = sortedLangs
	projecte.LOC = LinesOfCode
	projecte.doc = readmeFile.Path
	projecte.LastMod = dir.GetLastUpdate()
	projecte.Saved = saved
	projecte.Doc = make(map[string]string, 0)
	projecte.Path = dir.Path

	return projecte, nil
}
func (dir FileObject) Found(name string) (FileObject, error) {
	// Search for a specific file in a dir
	for _, item := range dir.Content {
		if item.Name == name {
			return item, nil
		}
	}

	return FileObject{}, errors.New("FILE NOT FOUND")
}

//  Real functions used outside
func IndexDirectory(dirPath string) []FileObject {
	list := []FileObject{}

	// Get every files in the directory requested
	files, err := os.ReadDir(dirPath)
	if err != nil {
		console.error("%s", err.Error())
	}

	for _, file := range files {
		// Ignore the automatic build dirs
		if slices.Contains(Configuration.Diskhub.Exclude, file.Name()) {
			continue
		}

		// Gather the information of the file
		fullPath := filepath.Join(dirPath, file.Name())
		infoFile, err := file.Info()
		if err != nil {
			console.fatal("%s", err.Error())
		}
		modTime := infoFile.ModTime().Local().Format("2006-01-02 15:04:05")

		// Add the file to the list
		if file.IsDir() {
			actual := FileObject{Name: file.Name(), Path: fullPath, IsDir: true, Content: IndexDirectory(fullPath), LastMod: modTime}
			list = append(list, actual)
		} else {
			list = append(list, FileObject{Name: file.Name(), Path: fullPath, IsDir: false, Content: []FileObject{}, LastMod: modTime})
		}
	}

	return list
}
func InitProjects() {
	// Loop on the projects global variables
	for _, path := range Configuration.Diskhub.Paths {
		dirs := IndexDirectory(path)
		for _, item := range dirs {
			// If it found the diskhub.toml file, add it, else ignore the dir
			projecte, err := item.EstablishProject()
			if err != nil {
				if err.Error() == "FILE NOT FOUND" { continue }
				console.error("%s", err.Error())
			} else {
				console.info("Project %s found in %s", projecte.Name, item.Path)
				projects = append(projects, projecte)
			}
		}
	}
}
