package usecases

import (
	"errors"
	"fmt"
	"github.com/mp-hl-2021/splinter/auth"
	"github.com/mp-hl-2021/splinter/storage"
	"golang.org/x/crypto/bcrypt"
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

type DummyUserInterface struct {
	Auth    auth.Interface
	Storage storage.Interface
}

func (d *DummyUserInterface) CreateAccount(username, password string) (User, error) {
	fmt.Printf("CreateAccount: %s %s\n\n", username, password)
	if err := validateUsername(username); err != nil {
		return User{}, err
	}
	if err := validatePassword(password); err != nil {
		return User{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	acc, err := d.Storage.CreateAccount(storage.Credentials{
		Username: username,
		Password: string(hashedPassword),
	})
	if err != nil {
		return User{}, err
	}
	return User{Id: UserId(acc.Id)}, nil
}

func (d *DummyUserInterface) Authenticate(username, password string) (Token, error) {
	fmt.Printf("Authenticate: %s %s\n\n", username, password)
	if err := validateUsername(username); err != nil {
		return "", err
	}
	if err := validatePassword(password); err != nil {
		return "", err
	}
	acc, err := d.Storage.GetAccountByUsername(username)
	if err != nil {
		return "", ErrInvalidLogin
	}
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Credentials.Password), []byte(password)); err != nil {
		return "", ErrInvalidPassword
	}
	token, err := d.Auth.IssueToken(acc.Id)
	return Token(token), err
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

func (d DummyUserInterface) Vote(snippet SnippetId, vote int) error {
	// TODO: implement me
	return nil
}

func (d DummyUserInterface) PostComment(contents string, snippet SnippetId) (Comment, error) {
	// TODO: implement me
	return Comment{Contents: contents, Snippet: snippet}, nil
}

func (d DummyUserInterface) GetComments(snippet SnippetId) ([]Comment, error) {
	// TODO: implement me
	return []Comment{}, nil
}

func (d DummyUserInterface) DeleteComment(comment CommentId) error {
	// TODO: implement me
	return nil
}
