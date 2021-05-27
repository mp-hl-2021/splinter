package usecases

import (
	"errors"
	"time"
)

type ProgrammingLanguage string
type UserId uint
type SnippetId uint
type CommentId uint
type Token string

var (
	ErrInvalidLogin    = errors.New("login not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type User struct {
	Id       UserId // Unique identifier, persists through username changes
	Username string // Visible username, can be changed
}

type Rating struct {
	Likes    int
	Dislikes int
}

type Comment struct {
	Id        CommentId
	Contents  string
	Snippet   SnippetId
	Author    UserId
	CreatedAt time.Time
}

type Snippet struct {
	Id              SnippetId
	Contents        string
	Language        ProgrammingLanguage
	Author          UserId
	Rating          Rating
	CurrentUserVote int
	CreatedAt       time.Time
}

type UserInterface interface {
	CreateAccount(username, password string) (User, error)
	Authenticate(username, password string) (Token, error)
	GetCurrentUser() (User, error)
	GetUser(user UserId) (User, error)

	PostSnippet(contents string, language ProgrammingLanguage) (Snippet, error)
	GetSnippetsByUser(user UserId) ([]Snippet, error)
	GetSnippetsByLanguage(language ProgrammingLanguage) ([]Snippet, error)
	GetSnippet(snippet SnippetId) (Snippet, error)
	DeleteSnippet(snippet SnippetId) error
	Vote(snippet SnippetId, vote int /* ±1 */) error

	PostComment(contents string, snippet SnippetId) (Comment, error)
	GetComments(snippet SnippetId) ([]Comment, error)
	DeleteComment(comment CommentId) error
}