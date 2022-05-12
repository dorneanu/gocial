package html

import (
	"embed"
	"errors"
	"html/template"
	"io"

	"github.com/dorneanu/gomation/internal/entity"
	"github.com/labstack/echo/v4"
	// "github.com/dorneanu/gomation/internal/jwt"
)

var (
	//go:embed templates/*
	Templates embed.FS

	//go:embed static/*
	StaticContent embed.FS

	// base template
	baseTemplate = "base.html"
)

type TemplateRegistry struct {
	templates map[string]*template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, baseTemplate, data)
}

type IndexParams struct {
	ProviderIndex entity.AuthProviderIndex
}

type SharePostParams struct {
	SendButtonMessage   string
	CancelButtonMessage string
}

func parse(file string) *template.Template {
	return template.Must(
		template.New(baseTemplate).ParseFS(Templates, "templates/"+baseTemplate, file))
}

// registerTemplates sets up html templating system
func RegisterTemplates() *TemplateRegistry {
	templates := make(map[string]*template.Template)
	templates["index"] = parse("templates/index.html")
	templates["about"] = parse("templates/about.html")
	templates["authInfo"] = parse("templates/auth/info.html")
	templates["shareIndex"] = parse("templates/share/index.html")

	return &TemplateRegistry{
		templates: templates,
	}
}
