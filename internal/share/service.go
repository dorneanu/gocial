package share

import (
	"context"

	"github.com/dorneanu/gomation/internal/entity"
)

type Service interface {
	ShareArticle(entity.ArticleShare) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{
		repo: repo,
	}
}

// ShareArticle shares an article using the specified repository
func (s service) ShareArticle(article entity.ArticleShare) error {
	return s.repo.ShareArticle(context.Background(), article)
}
