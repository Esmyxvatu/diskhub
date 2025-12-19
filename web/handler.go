package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/esmyxvatu/feather"
	"github.com/yuin/goldmark"
	mdHighlight "github.com/yuin/goldmark-highlighting/v2"
)

//============================================================ Global Variables =====================================================================
var Ram_Values 	  []StringInt
var Max_RAM_Value int		  = 50

//============================================================ Type definitions =====================================================================

type IndexLangs struct {
	Name  	string
	Loc   	int
	Files 	int
	Projs 	int
}
type IndexPage struct {
	Projects 	[]Project
	Info    	map[string]any
	Langs    	[]IndexLangs
	StatusChart template.HTML
	Status      map[string]float64
	Cookies		map[string]any

}
type SearchPage struct {
	SearchTerm string
	Projects   []Project
	Time       int64
	Cookies    map[string]any

}
type ProjectPage struct {
	Name    string
	Id      string
	About   string
	Links   map[string]string
	Files   []FileObject
	Langs   []StringInt
	Tags    []string
	LOC     int
	Saved   bool
	Readme  template.HTML
	Path    string
	Status  string
	Cookies map[string]any

}
type SingleFilePage struct {
	Name     string
	Id       string
	Files    []FileObject
	Path     string
	Content  template.HTML
	FileName string
	Info     map[string]any
	Cookies  map[string]any
}
type DocPage struct {
	Articles []string
	Project  Project
	Name     string
	Content  template.HTML
	Cookies  map[string]any

}
type FileExp struct {
	Name 		string
	Path 		string
	IsDir 		bool
	Modified	string
	Size 		int64
}
type FileExpPage struct {
	Explorer []FileExp
	Cookies  map[string]any
}
type DashboardPage struct {
	Uptime  time.Duration
	RAM     float32
	Cookies map[string]any
	Chart   template.HTML
}
type SnippetPage struct {
	Cookies			map[string]any
	SearchTerm  	string
	IdSelected  	string
	SnippetChoosed  Snippet
	SnippetFound    []Snippet
	SnippetContent  template.HTML
}
type ChatPage struct {
	Models		[]string
	Messages	[]Message
	Cookies		map[string]any
}

//============================================================ User Routes ==========================================================================

// Method accessible for the user
func IndexHandler(ctx *feather.Context) {
	// Calculate the different information
	var nbfiles int
	var nbloc int
	for i := range langStat {
		nbfiles += langStat[i].Files
		nbloc += langStat[i].Loc
	}
	infos := map[string]any{
		"nbLangs": len(langStat),
		"nbFiles": nbfiles,
		"nbLOC":   nbloc,
	}

	// Get the status of every projects, and save it in a dict
	statusList := map[string]float64{}
	found := []string{}
	for _, proj := range projects {
		if slices.Contains(found, proj.Status) {
			statusList[proj.Status]++;
		} else {
			statusList[proj.Status] = 1
			found = append(found, proj.Status)
		}
	}

	// Check for the theme cookie
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err == nil { if themeCookie.Value == "dark" { dark = true } }

	// Create the cookies dict
	cookies := map[string]any{
		"Theme": dark,
	}

	// Create the dict
	data := IndexPage{
		Projects:	 projects,
		Info:    	 infos,
		Langs:    	 langStat,
		StatusChart: template.HTML(generatePieChart(statusList, 300, 300)),		// Escape the generated pie with template.HTML to be interpreted as HTML
		Status:   	 statusList,
		Cookies:     cookies,
	}

	// Send the template
	ctx.Template([]string{"website/index.html"}, data, nil)
}
func SearchHandler(ctx *feather.Context) {
	// Get results and mesures the time taken for the search
	start := time.Now().UnixMilli()
	results, err := SearchFor(ctx.Query("q"))
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}
	end := time.Now().UnixMilli()

	// Load the theme cookie
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err == nil { if themeCookie.Value == "dark" { dark = true } }

	// Create the map of the Cookies
	cookies := map[string]any{
		"Theme": dark,
	}

	// Unify all data in one variable
	data := SearchPage{
		SearchTerm: ctx.Query("q"),
		Projects:   results,
		Time:       end - start,
		Cookies:    cookies,
	}

	// Send the template
	ctx.Template([]string{"website/search.html"}, data, nil)
}
func ProjectHandler(ctx *feather.Context) {
	// Search the requested project
	var projecte Project
	found := false
	for _, proj := range projects {
		if proj.Id == ctx.Params["id"] {
			projecte = proj
			found = true
			break
		}
	}
	if !found {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	// Load the ReadMe content from the file
	content, err := os.ReadFile(projecte.doc)
	if err != nil {
		console.error("Error reading file: %s", err.Error())
		content = []byte{42, 42, 78, 111, 32, 82, 101, 97, 100, 109, 101, 32, 70, 111, 117, 110, 100, 42, 42}
	}

	// Some code for the template
	projecte.Langs[len(projecte.Langs)-1].IsLast = true

	// Load the theme cookie
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err == nil { if themeCookie.Value == "dark" { dark = true } }

	// Create the cookie dict
	cookies := map[string]any{
		"Theme": dark,
	}

	// Unify the data
	data := ProjectPage{
		Name:    projecte.Name,
		Id:      projecte.Id,
		About:   projecte.About,
		Links:   projecte.Links,
		Files:   projecte.Files,
		Langs:   projecte.Langs,
		LOC:     projecte.LOC,
		Saved:   projecte.Saved,
		Tags:    projecte.Tags,
		Readme:  markdownToHTML(content),
		Path:    projecte.Path,
		Status:  projecte.Status,
		Cookies: cookies,
	}

	// Send the data
	ctx.Template([]string{"website/project.html"}, data, nil)
}
func FileShowerHandler(ctx *feather.Context) {
	// Search for the project
	var projecte Project
	found := false
	for _, proj := range projects {
		if proj.Id == ctx.Params["id"] {
			projecte = proj
			found = true
			break
		}
	}
	if !found {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	// Search for the file
	var file FileObject
	found = false
	parts := strings.Split(ctx.Params["filepath"], "/")		// Split the path | "/abc/114/e.html" => [abc 114 e.html]
	if len(parts) > 1 {
		files := projecte.Files

		for _, part := range parts {			// Search for every parts
			for _, fily := range files {		// Search for every file
				if fily.Name == part {			// Try if it's the good file
					if fily.IsDir {				// If it's a dir, return the list of files
						files = fily.Content
					} else {					// If it's a file, return it
						file = fily
						found = true
						break
					}
				}
			}
			if found {
				break
			}
		}

		// Throw an error to the client
		if !found {
			http.NotFound(ctx.Writer, ctx.Request)
			return
		}
	} else {
		// If it's a single file ("a.html") just check in the proj list
		for _, fily := range projecte.Files {
			if fily.Name == ctx.Params["filepath"] {
				file = fily
				found = true
				break
			}
		}

		// Throw an error if not found
		if !found {
			http.NotFound(ctx.Writer, ctx.Request)
			return
		}
	}

	// Read the file
	content, err := os.ReadFile(file.Path)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
	}
	
	// Gather the information relative to the file
	langName, loc := file.GetLangInfo()
	information := map[string]any{
		"lines": len(strings.Split(string(content), "\n")),
		"LOC":   loc,
		"lang":  langName,
	}

	// Load the theme cookie
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err == nil { if themeCookie.Value == "dark" { dark = true } }

	// Create the cookies dict
	cookies := map[string]any{
		"Theme": dark,
	}

	// Unify data
	data := SingleFilePage{
		Name:     projecte.Name,
		Id:       projecte.Id,
		Files:    projecte.Files,
		Path:     file.Path,
		Content:  stringToChroma(string(content), file.Name),
		FileName: file.Name,
		Info:     information,
		Cookies:  cookies,
	}

	// Send the template
	ctx.Template([]string{"website/singlefile.html", "website/dirfileshower.html"}, data, nil)
}
func DocShowerHandler(ctx *feather.Context) {
	// Search for the project
	var projecte Project
	found := false
	for _, proj := range projects {
		if proj.Id == ctx.Params["id"] {
			projecte = proj
			found = true
			break
		}
	}
	if !found {
		console.error("Project not found for ID: %s", ctx.Params["id"])
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	// Get the list of all articles of the doc and the content of the one asked
	articles := []string{}
	var contentMd template.HTML = template.HTML("<pre>File not found</pre>")
	for article := range projecte.Doc {
		articles = append(articles, article)

		// If it's the asked one
		if article == ctx.Params["page"] {
			// Search for the file
			var file FileObject
			found = false
			parts := strings.Split(projecte.Doc[article], "/")
			if len(parts) > 1 {
				files := projecte.Files

				// Search to get the right file
				for _, part := range parts {
					for _, fily := range files {
						if fily.Name == part {
							if fily.IsDir {
								files = fily.Content
							} else {
								file = fily
								found = true
								break
							}
						}
					}
					if found {
						break
					}
				}

				// Throw an error if the file isn't found
				if !found {
					console.error("File not found: %s (Article: %s, Associated file: %s)", file.Path, article, projecte.Doc[article])
					http.NotFound(ctx.Writer, ctx.Request)
					return
				}
			} else {
				for _, fily := range projecte.Files {
					if fily.Name == projecte.Doc[article] {
						file = fily
						found = true
						break
					}
				}

				// Throw an error if not found
				if !found {
					console.error("File not found: %s (Article: %s, Associated file: %s)", file.Path, article, projecte.Doc[article])
					http.NotFound(ctx.Writer, ctx.Request)
					return
				}
			}

			// Read the content
			content, err := os.ReadFile(file.Path)
			if err != nil {
				ctx.Error(http.StatusInternalServerError, err.Error())
			}

			// Convert markdown to HTML or load the file as HTML
			if strings.HasSuffix(file.Name, ".md") {
				contentMd = markdownToHTML(content)
			} else {
				contentMd = template.HTML("<pre>" + string(content) + "</pre>")
			}
			ctx.Params["page"] = article
		}
	}

	// Load the theme cookie
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err == nil { if themeCookie.Value == "dark" { dark = true } }

	// Create the cookies dict
	cookies := map[string]any{
		"Theme": dark,
	}

	// Unify all the data
	data := DocPage{
		Project:  projecte,
		Articles: articles,
		Name:     ctx.Params["page"],
		Content:  contentMd,
		Cookies:  cookies,
	}

	// Send the template
	ctx.Template([]string{"website/docs.html"}, data, nil)
}
func FileExplorerHandler(ctx *feather.Context) {
	// Init the various needed variable
	path := strings.ReplaceAll(ctx.Params["path"], "%20", " ")
	
	console.info("Param: %s, final path: %s", ctx.Params["path"], path)

	// Define the way to go back
	explorer := []FileExp{
		{
			Name:     "..",
			Path:     strings.Join(strings.Split(path, "/")[0 : len(strings.Split(path, "/")) - 1], "/"),
			IsDir:    true,
			Modified: time.Now().Format(time.ANSIC),
			Size:     0,
		},
	}

	// Read the content of the repertory
	files, err := os.ReadDir(path)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}

		// Add the file to the list
		if file.IsDir() {
			explorer = append(explorer, FileExp{
				Name: file.Name(),
				Path: path + "/" + file.Name(),
				IsDir: true,
				Modified: fileInfo.ModTime().Format(time.ANSIC),
				Size: 0,
			})
		} else {
			explorer = append(explorer, FileExp{
				Name: file.Name(),
				Path: path + "/" + file.Name(),
				IsDir: false,
				Modified: fileInfo.ModTime().Format(time.ANSIC),
				Size: fileInfo.Size(),
			})
		}
	}

	// Check for the theme cookie
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err == nil { if themeCookie.Value == "dark" { dark = true } }

	// Create the cookies dict
	cookies := map[string]any{
		"Theme": dark,
	}

	// Unify the data
	data := FileExpPage{
		Explorer: explorer,
		Cookies:  cookies,
	}

	// Send the template
	ctx.Template([]string{"website/fileexplorer.html"}, data, template.FuncMap{
		"processSize": processSize,
	})
}
func SnippetHandler(ctx *feather.Context) {
	// Check for the theme cookie
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err == nil { if themeCookie.Value == "dark" { dark = true } }

	// Create the cookies dict
	cookies := map[string]any{
		"Theme": dark,
	}

	choosed := GetSnippet(ctx.Query("id"))

	data := SnippetPage{
		Cookies: 		cookies,
		SearchTerm: 	ctx.Query("q"),
		IdSelected: 	ctx.Query("id"),
		SnippetChoosed: choosed,
		SnippetFound: 	GetAllSnippet(),
		SnippetContent: stringToChroma(choosed.Content, "lol.js"),
	}

	// Send the template
	ctx.Template([]string{"website/snippet.html"}, data, nil)
}
func ChatHandler(ctx *feather.Context) {
	// Check for the theme cookie
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err == nil { if themeCookie.Value == "dark" { dark = true } }

	// Create the cookies dict
	cookies := map[string]any{
		"Theme": dark,
	}

	// Unify the data
	data := ChatPage{
		Models:		ListModelAvaillable(),
		Messages:   Messages,
		Cookies:    cookies,
	}

	// Send the template
	ctx.Template([]string{"website/chat.html"}, data, nil)
}

//============================================================ API Routes ===========================================================================

func APISaveHandler(ctx *feather.Context) {
	id := ctx.FormValue("id")

	// Modify the project
	for i := range projects {
		if projects[i].Id == id {
			projects[i].Saved = !projects[i].Saved

			file, err := os.OpenFile(projects[i].Path + "/.diskhub/.id", os.O_WRONLY, os.ModeAppend)
			console.verify(err)
			defer file.Close()

			_, err = file.WriteString(id + "\ntrue")
			console.verify(err)

			break
		}
	}

	// Redirect the user
	ctx.Redirect(http.StatusSeeOther, "/project/"+id)
}
func APIRefreshHandler(ctx *feather.Context) {
	// Reinitialize the global variables
	Langs = []Language{}
	projects = []Project{}
	langStat = []IndexLangs{}

	// Reload the different global variables
	InitLangs()
	InitProjects()
	InitStats()

	// Redirect the user
	ctx.Redirect(http.StatusSeeOther, "/")
}
func APISetCookie(ctx *feather.Context) {
	//==================== Theme ====================================
	theme := ctx.Query("theme")

	if theme != "dark" && theme != "light" {
		ctx.Error(http.StatusBadRequest, "Invalid theme")
		return
	}

	// Define the cookie
	cookie := &http.Cookie{
		Name: 	"theme",
		Value:	theme,
		Path:  	"/",
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}

	// Set the cookie and redirect the user
	ctx.SetCookie(cookie)
	ctx.Redirect(http.StatusSeeOther, "/")
}
func APIAskOllama(ctx *feather.Context) {
	var data UserQuestion
	ctx.JSONBody(data)

	answer := AskOllama(data.Model, data.Content, data.Stream)

	ctx.JSON(http.StatusOK, answer)
}

//============================================================ Admin Routes =========================================================================

func AdminDashboardHandler(ctx *feather.Context) {
	// Load the theme cookie
	dark := false
	themeCookie, err := ctx.Cookie("theme")
	if err == nil { if themeCookie.Value == "dark" { dark = true } }

	// Create the cookies dict
	cookies := map[string]any{
		"Theme": dark,
	}

	// Load and read the stats for the usage of memory
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	Ram_Values = append(Ram_Values, StringInt{"", float32(m.Alloc/1024), false})
	if len(Ram_Values) > Max_RAM_Value { Ram_Values = slices.Delete(Ram_Values, 0, 1) }

	// Unify all the data
	data := DashboardPage{
		Uptime:  time.Since(Uptime_start),
		RAM:	 Ram_Values[len(Ram_Values)-1].Value,
		Cookies: cookies,
		Chart:   template.HTML(generateLineChart(Ram_Values, 480, 200)),
	}

	// Send the template
	ctx.Template([]string{"website/dashboard.html"}, data, nil)
}

//============================================================ Helper Functions =====================================================================
func processSize(size int64) string {
	// Check to convert the bytes to more humand readable sizes (89128960 B => 85 MB)
	if size < 1024 {
		return fmt.Sprint(size) + " B"
	} else if size < 1024*1024 {
		return fmt.Sprint(size/1024) + " KB"
	} else if size < 1024*1024*1024 {
		return fmt.Sprint(size/(1024*1024)) + " MB"
	} else {
		return fmt.Sprint(size/(1024*1024*1024)) + " GB"
	}
}
func markdownToHTML(mdContent []byte) template.HTML {
	// Convert markdown to HTML
	// htmlContent := blackfriday.Run(mdContent)

	markdown := goldmark.New(
		goldmark.WithExtensions( mdHighlight.Highlighting ),
	)
	var byteContent bytes.Buffer
	if err := markdown.Convert(mdContent, &byteContent); err != nil { console.fatal("%v", err) }
	htmlContent := byteContent.String()

	return template.HTML(htmlContent)
}
func stringToChroma(content string, name string) template.HTML {
	// Define the lexer
	lexer := lexers.Match(name)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Load a theme style
	style := styles.Get("swapoff")
	if style == nil {
		style = styles.Fallback
	}

	// Define the formater to uses classes
	formatter := html.New(
		html.WithClasses(true),
		html.WithAllClasses(true),
	)

	// Tokenize the content
	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		console.fatal("%s", err.Error())
	}

	// Format the tokens
	var w bytes.Buffer
	err = formatter.Format(&w, style, iterator)
	if err != nil {
		console.fatal("%s", err.Error())
	}

	// Convert to HTML
	return template.HTML(w.String())
}
