package handlers

import (
	"html/template"
	"path/filepath"

	"github.com/snirkop89/shopping-app/pkg/repository"
)

var tmpl *template.Template

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{Repo: repo}
}

func init() {
	templatesDir := "./templates"
	pattern := filepath.Join(templatesDir, "**", "*.html")
	tmpl = template.Must(template.ParseGlob(pattern))
}
