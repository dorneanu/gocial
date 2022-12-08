package share

import (
	"context"

	"github.com/dorneanu/gocial/internal/entity"
)

type Repository interface {
	ShareArticle(context.Context, entity.ArticleShare) error
}
