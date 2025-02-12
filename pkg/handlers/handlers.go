package handlers

import "github.com/snirkop89/shopping-app/pkg/repository"

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{Repo: repo}
}
