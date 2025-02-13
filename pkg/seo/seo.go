package seo

import (
	"html/template"
)

type GlobalSEO struct {
	GlobalTitle  string          `json:"title"`
	GlobalExtras []template.HTML `json:"global"`
}

type PageSEO struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Keywords    []string        `json:"keywords"`
	Extra       []template.HTML `json:"extra"`
}

type SEO struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Keywords    []string        `json:"keywords"`
	Extra       []template.HTML `json:"extra"`
}

func CombineSeo(g GlobalSEO, p PageSEO) SEO {
	extras := g.GlobalExtras
	extras = append(extras, p.Extra...)

	var combineTitle string
	if len(p.Title) != 0 {
		combineTitle = g.GlobalTitle + " - " + p.Title
	} else {
		combineTitle = g.GlobalTitle
	}

	s := SEO{
		Title:       combineTitle,
		Description: p.Description,
		Keywords:    p.Keywords,
		Extra:       extras,
	}

	return s
}
