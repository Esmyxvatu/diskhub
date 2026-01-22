package wiki

import (
	"diskhub/web/models"
	"diskhub/web/render"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
)

func Generate(dir models.FileObject) ([]any, map[string]*models.DocArticle, error) {
	articles := []any{}
	groups := []any{}
	pathMap := map[string]*models.DocArticle{}

	wikiDir, err := dir.Found("wiki")
	if err != nil {
		return []any{}, make(map[string]*models.DocArticle), err
	}

	if !wikiDir.IsDir {
		return []any{}, make(map[string]*models.DocArticle), errors.New("'wiki' should be a folder and not a file in .diskhub")
	}

	for _, object := range wikiDir.Content {
		if object.IsDir {
			group := models.DocGroup{IsArticle: false}
			group.Name = strings.Title(strings.Join(strings.Split(object.Name, "_"), " "))
			group_prefix := strings.ToLower(strings.Join(strings.Split(group.Name, " "), "_"))

			for _, subObject := range object.Content {
				if subObject.IsDir {
					continue
				}

				article, err := CreateArticle(subObject)
				if err != nil {
					return []any{}, make(map[string]*models.DocArticle), err
				}

				article.Path = fmt.Sprintf("%s/%s", group_prefix, article.Path)
				pathMap[article.Path] = &article

				group.Articles = append(group.Articles, article)
			}

			groups = append(groups, group)
		} else {
			article, err := CreateArticle(object)
			if err != nil {
				return []any{}, make(map[string]*models.DocArticle), err
			}

			articles = append(articles, article)

			pathMap[article.Path] = &article
		}
	}

	// Tri docs pour avoir les articles au dessus
	docs := []any{}
	docs = append(docs, articles...)
	docs = append(docs, groups...)
	OrganizeArticles(&docs, wikiDir)

	return docs, pathMap, nil
}

func CreateArticle(object models.FileObject) (models.DocArticle, error) {
	article := models.DocArticle{IsArticle: true}
	name, extension, _ := strings.Cut(object.Name, ".")
	article.Name = strings.Title(strings.Join(strings.Split(name, "_"), " "))
	article.Path = strings.ToLower(name)

	bytes, err := os.ReadFile(object.Path)
	if err != nil {
		return models.DocArticle{}, err
	}

	if extension == "md" {
		article.Content = render.MarkdownToHTML(bytes)
	} else {
		article.Content = template.HTML(bytes)
	}

	return article, nil
}
