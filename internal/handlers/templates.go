package handlers

import (
	"html/template"
	"io"
	"github.com/labstack/echo/v4"
)

// TemplateRenderer handles rendering of HTML templates
type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplates() *TemplateRenderer {
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

	return &TemplateRenderer{
		templates: tmpl,
	}
}