package share

import (
	"context"

	"github.com/dorneanu/gomation/internal/entity"
)

type Service interface {
	ShareArticle(entity.ArticleShare) error
	ShareComment(entity.CommentShare) error
}

type shareService struct {
	repo Repository
}

func NewShareService(repo Repository) Service {
	return shareService{
		repo: repo,
	}
}

// ShareArticle shares an article using the specified repository
func (s shareService) ShareArticle(article entity.ArticleShare) error {
	return s.repo.ShareArticle(context.Background(), article)
}

// TODO: Implement ShareComment ...
func (s shareService) ShareComment(comment entity.CommentShare) error {
	return nil
}
