package types

type SnippetStorage interface {
	AddSnippet(snippet Snippet) (SnippetId, error)
	GetSnippetsByUser(user UserId) ([]Snippet, error)
	GetSnippetsByLanguage(language ProgrammingLanguage) ([]Snippet, error)
	GetSnippet(snippet SnippetId) (Snippet, error)
	DeleteSnippet(snippet SnippetId) error
	SetSnippetHighlight(snippet SnippetId, highlight string) error
	Vote(user UserId, snippet SnippetId, vote int) error
	GetVote(user UserId, snippet SnippetId) (int, error)

	AddComment(comment Comment) (CommentId, error)
	GetComment(comment CommentId) (Comment, error)
	GetComments(snippet SnippetId) ([]Comment, error)
	DeleteComment(comment CommentId) error
}
