package entity

// ArticleShare is an article to be shared via the share service
type ArticleShare struct {
	URL     string
	Title   string
	Comment string
}

// CommentShare is a comment to be shared via the share service
type CommentShare struct {
	// TODO: Any other fields needed?
	Comment string
}
