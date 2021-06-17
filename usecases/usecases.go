package usecases

import "github.com/mp-hl-2021/splinter/types"

type UserInterface interface {
	CreateAccount(username, password string) (types.User, error)
	Authenticate(username, password string) (types.Token, error)
	GetUser(user types.UserId) (types.User, error)

	PostSnippet(author types.UserId, contents string, language types.ProgrammingLanguage) (types.Snippet, error)
	GetSnippetsByUser(user types.UserId, current types.UserId) ([]types.Snippet, error)
	GetSnippetsByLanguage(language types.ProgrammingLanguage, current types.UserId) ([]types.Snippet, error)
	GetSnippet(current types.UserId, snippet types.SnippetId) (types.Snippet, error)
	DeleteSnippet(current types.UserId, snippet types.SnippetId) error
	Vote(current types.UserId, snippet types.SnippetId, vote int /* ±1 */) error

	PostComment(author types.UserId, contents string, snippet types.SnippetId) (types.Comment, error)
	GetComments(snippet types.SnippetId) ([]types.Comment, error)
	DeleteComment(current types.UserId, comment types.CommentId) error
}
