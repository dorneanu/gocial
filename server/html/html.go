package html

import (
	"embed"
	"html/template"
	"io"

	"github.com/dorneanu/gomation/internal/entity"
	"github.com/dorneanu/gomation/internal/jwt"
	// "github.com/dorneanu/gomation/internal/jwt"
)

var (
	//go:embed templates/*
	Templates embed.FS

	//go:embed static/*
	StaticContent embed.FS

	profile = parse("templates/profile.html")
	index   = parse("templates/index.html")
	post    = parse("templates/post.html")
)

type IndexParams struct {
	ProviderIndex entity.AuthProviderIndex
}

type PostParams struct {
	SendButtonMessage   string
	CancelButtonMessage string
}

func Index(w io.Writer, p IndexParams) error {
	return index.Execute(w, p)
}

func Profile(w io.Writer, jwtClaims jwt.JwtCustomClaims) error {
	return profile.Execute(w, jwtClaims)
}

func ArticlePost(w io.Writer, p PostParams) error {
	return post.Execute(w, p)
}

func parse(file string) *template.Template {
	return template.Must(
		template.New("layout.html").ParseFS(Templates, "templates/layout.html", file))
}
