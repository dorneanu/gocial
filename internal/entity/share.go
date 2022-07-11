package entity

// ArticleShare is an article to be shared via the share service
type ArticleShare struct {
	URL       string `json:"url" form:"url" validate:"required"`
	Title     string `json:"title" form:"title" validate:"required"`
	Comment   string `json:"comment" form:"comment" validate:"required"`
	Providers string `json:"providers" form:"providers" validate:"required"`
}

// CommentShare is a comment to be shared via the share service
type CommentShare struct {
	// TODO: Any other fields needed?
	Comment string
}
