package handler

import (
	"diskhub/web/language"
	"diskhub/web/render"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/esmyxvatu/feather"
)

type FileExp struct {
	Name     string
	Path     string
	IsDir    bool
	Modified string
	Size     int64
}
type RealFileExp struct {
	Name        string
	Perms       string
	Modified    string
	Size        int64
	ContentType string
	Content     template.HTML
	Back        string
}
type FileExpPage struct {
	IsExplorer bool
	FileToShow RealFileExp
	Explorer   []FileExp
	Cookies    map[string]any
}

func FileExplorerHandler(ctx *feather.Context) {
	// Init the various needed variable
	path := strings.ReplaceAll(ctx.Params["path"], "%20", " ")
	if path == "" {
		path = "/"
	}

	// Define the way to go back
	explorer := []FileExp{
		{
			Name:     "..",
			Path:     strings.Join(strings.Split(path, "/")[0:len(strings.Split(path, "/"))-1], "/"),
			IsDir:    true,
			Modified: time.Now().Format(time.ANSIC),
			Size:     0,
		},
	}
	finalFile := RealFileExp{}

	info, err := os.Stat(path)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	// Read the content of the repertory
	if info.IsDir() {
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
					Name:     file.Name(),
					Path:     path + "/" + file.Name(),
					IsDir:    true,
					Modified: fileInfo.ModTime().Format(time.ANSIC),
					Size:     0,
				})
			} else {
				explorer = append(explorer, FileExp{
					Name:     file.Name(),
					Path:     path + "/" + file.Name(),
					IsDir:    false,
					Modified: fileInfo.ModTime().Format(time.ANSIC),
					Size:     fileInfo.Size(),
				})
			}
		}
	} else {
		file, err := os.Open(path)
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}

		fileInfo, err := file.Stat()
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}

		fileSize := fileInfo.Size()
		fileContent := make([]byte, fileSize)
		_, err = file.Read(fileContent)
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}

		fileMIME, _, _ := strings.Cut(http.DetectContentType(fileContent), ";")
		content := render.StringToChroma(string(fileContent))
		if strings.HasPrefix(fileMIME, "text") {
			ext := strings.TrimPrefix(filepath.Ext(file.Name()), ".")
			lang := language.LangByExt[ext]

			if lang != nil && lang.Name == "Markdown" {
				content = render.MarkdownToHTML(fileContent)
			}
		}

		if fileMIME == "application/octet-stream" {
			content = render.DumpToHTML(fileContent)
		}

		finalFile = RealFileExp{
			Name:        fileInfo.Name(),
			Perms:       fileInfo.Mode().Perm().String(),
			Modified:    fileInfo.ModTime().Format(time.ANSIC),
			Size:        fileSize,
			ContentType: fileMIME,
			Content:     content,
			Back:        strings.Join(strings.Split(path, "/")[0:len(strings.Split(path, "/"))-1], "/"),
		}
	}

	cookies, err := makeCookieMap(ctx)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	// Unify the data
	data := FileExpPage{
		IsExplorer: info.IsDir(),
		FileToShow: finalFile,
		Explorer:   explorer,
		Cookies:    cookies,
	}

	// Send the template
	ctx.Template([]string{"website/fileexplorer.html"}, data, template.FuncMap{
		"processSize": processSize,
	})
}

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
