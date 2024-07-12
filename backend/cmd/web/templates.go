package main

import (
	"Projet-Forum/internal/data"
	"Projet-Forum/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

var functions = template.FuncMap{
	"humanDate":       humanDate,
	"getUserReaction": getUserReaction,
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func getUserReaction(user data.User, postID int) string {
	if user.Reactions != nil {
		emoji, ok := user.Reactions[postID]
		if ok {
			return emoji
		}
	}
	return ""
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "templates/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"templates/base.tmpl",
			"templates/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
