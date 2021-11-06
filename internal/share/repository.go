package share

import (
	"context"

	"github.com/dorneanu/gomation/internal/entity"
)

type Repository interface {
	ShareArticle(context.Context, entity.ArticleShare) error
}
