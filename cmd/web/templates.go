package main

import (
	"html/template"
	"path/filepath"

	"github.com/rostamborn/snippetbox/pkg/forms"
	"github.com/rostamborn/snippetbox/pkg/models"
)

type templateData struct {
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	CurrentYear       int
	Form              *forms.Form
	Flash             string
	AuthenticatedUser *models.User
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
