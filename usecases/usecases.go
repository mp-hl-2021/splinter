package usecases

import "time"

type ProgrammingLanguage string
type UserId string
type SnippetId string
type CommentId string
type Token string

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
	Vote(snippet Snippet, vote int /* ±1 */) error

	PostComment(contents string, snippet SnippetId) (Comment, error)
	DeleteComment(comment CommentId) error
}

type DummyUserInterface struct{}

func (d DummyUserInterface) CreateAccount(username, password string) (User, error) {
	// TODO: implement me
	return User{Username: username}, nil
}

func (d DummyUserInterface) Authenticate(username, password string) (Token, error) {
	// TODO: implement me
	return "", nil
}

func (d DummyUserInterface) GetCurrentUser() (User, error) {
	// TODO: implement me
	return User{Username: "anonymous"}, nil
}

func (d DummyUserInterface) GetUser(user UserId) (User, error) {
	// TODO: implement me
	return User{Id: user, Username: "anonymous"}, nil
}

func (d DummyUserInterface) PostSnippet(contents string, language ProgrammingLanguage) (Snippet, error) {
	// TODO: implement me
	return Snippet{Contents: contents, Language: language}, nil
}

func (d DummyUserInterface) GetSnippetsByUser(user UserId) ([]Snippet, error) {
	// TODO: implement me
	return []Snippet{}, nil
}

func (d DummyUserInterface) GetSnippetsByLanguage(language ProgrammingLanguage) ([]Snippet, error) {
	// TODO: implement me
	return []Snippet{}, nil
}

func (d DummyUserInterface) GetSnippet(snippet SnippetId) (Snippet, error) {
	// TODO: implement me
	return Snippet{}, nil
}

func (d DummyUserInterface) DeleteSnippet(snippet SnippetId) error {
	// TODO: implement me
	return nil
}

func (d DummyUserInterface) Vote(snippet Snippet, vote int) error {
	// TODO: implement me
	return nil
}

func (d DummyUserInterface) PostComment(contents string, snippet SnippetId) (Comment, error) {
	// TODO: implement me
	return Comment{Contents: contents, Snippet: snippet}, nil
}

func (d DummyUserInterface) DeleteComment(comment CommentId) error {
	// TODO: implement me
	return nil
}
