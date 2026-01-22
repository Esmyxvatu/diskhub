package models

import (
	"html/template"
)

type DocArticle struct {
	Path      string
	Name      string
	Content   template.HTML
	IsArticle bool
}

type DocGroup struct {
	Name      string
	Articles  []DocArticle
	IsArticle bool
}

type Wiki struct {
	Articles  []any
	PathMap   map[string]*DocArticle
	OriginDir FileObject
}
