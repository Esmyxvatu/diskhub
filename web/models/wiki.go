package models

import (
	"html/template"
)

type DocElement interface {
	GetName() string
}

type DocArticle struct {
	Path      string
	Name      string
	Content   template.HTML
	IsArticle bool
}

func (article DocArticle) GetName() string { return article.Name }

type DocGroup struct {
	Name      string
	Articles  []DocArticle
	IsArticle bool
}

func (group DocGroup) GetName() string { return group.Name }

type Wiki struct {
	Articles  []DocElement
	PathMap   map[string]*DocArticle
	OriginDir FileObject
}
