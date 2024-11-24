package handlers

import (
	"fmt"
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

func NewTemplates() *TemplateRenderer {
	// Create a new template instance
	tmpl := template.New("").Funcs(template.FuncMap{
		"contains": func(slice []string, item string) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
        "safehtml": func(s string) template.HTML {
            return template.HTML(s)
        },
	})

	// Parse all templates in all subdirectories
	patterns := []string{
		"views/layouts/*.html",  // Base templates
		"views/pages/*.html",    // Full pages
		"views/partials/*.html", // HTMX partial components
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
