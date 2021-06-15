package types

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
	Id                  SnippetId
	Contents            string
	HighlightedContents string
	Language            ProgrammingLanguage
	Author              UserId
	Rating              Rating
	CurrentUserVote     int
	CreatedAt           time.Time
}

