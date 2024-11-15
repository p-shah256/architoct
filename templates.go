// templates.go
package main

import (
	"html/template"
	"io"
	"github.com/labstack/echo/v4"
)

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplates() *Templates {
	// Create a new template instance
	tmpl := template.New("")

	// Parse all templates in all subdirectories
	patterns := []string{
		"views/layouts/*.html",    // Base templates
		"views/pages/*.html",     // Full pages
		"views/partials/*.html",  // HTMX partial components
	}

	for _, pattern := range patterns {
		_, err := tmpl.ParseGlob(pattern)
		if err != nil {
			panic(err)
		}
	}

	return &Templates{
		templates: tmpl,
	}
}
