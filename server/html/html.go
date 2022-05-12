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

	profile   = parse("templates/profile.html")
	index     = parse("templates/index.html")
	post      = parse("templates/post.html")
	AboutPage = parse("templates/about.html")
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
	return tmpl.ExecuteTemplate(w, "layout.html", data)
}

type IndexParams struct {
	ProviderIndex entity.AuthProviderIndex
}

type PostParams struct {
	SendButtonMessage   string
	CancelButtonMessage string
}

func parse(file string) *template.Template {
	return template.Must(
		template.New("layout.html").ParseFS(Templates, "templates/layout.html", file))
}

// registerTemplates sets up html templating system
func RegisterTemplates() *TemplateRegistry {
	templates := make(map[string]*template.Template)
	templates["index"] = parse("templates/index.html")
	templates["about"] = parse("templates/about.html")
	templates["profile"] = parse("templates/profile.html")
	templates["post"] = parse("templates/post.html")

	return &TemplateRegistry{
		templates: templates,
	}
}
