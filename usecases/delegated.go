package usecases

import (
	"errors"
	"github.com/mp-hl-2021/splinter/auth"
	"golang.org/x/crypto/bcrypt"
)

var (
	MustBeSnippetAuthorErr = errors.New("must be snippet's author")
	MustBeCommentAuthorErr = errors.New("must be comment's author")
)

type SnippetStorage interface {
	AddSnippet(snippet Snippet) (SnippetId, error)
	GetSnippetsByUser(user UserId) ([]Snippet, error)
	GetSnippetsByLanguage(language ProgrammingLanguage) ([]Snippet, error)
	GetSnippet(snippet SnippetId) (Snippet, error)
	DeleteSnippet(snippet SnippetId) error
	Vote(user UserId, snippet SnippetId, vote int) error
	GetVote(user UserId, snippet SnippetId) (int, error)

	AddComment(comment Comment) (CommentId, error)
	GetComment(comment CommentId) (Comment, error)
	GetComments(snippet SnippetId) ([]Comment, error)
	DeleteComment(comment CommentId) error
}

type DelegatedUserInterface struct {
	Auth           auth.Authenticator
	UserStorage    auth.UserStorage
	SnippetStorage SnippetStorage
}

func (u *DelegatedUserInterface) CreateAccount(username, password string) (User, error) {
	if err := auth.ValidateUsername(username); err != nil {
		return User{}, err
	}
	if err := auth.ValidatePassword(password); err != nil {
		return User{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	acc, err := u.UserStorage.CreateAccount(auth.Credentials{
		Username: username,
		Password: string(hashedPassword),
	})
	if err != nil {
		return User{}, err
	}
	return User{Id: UserId(acc.Id), Username: username}, nil
}

func (u *DelegatedUserInterface) Authenticate(username, password string) (Token, error) {
	if err := auth.ValidateUsername(username); err != nil {
		return "", err
	}
	if err := auth.ValidatePassword(password); err != nil {
		return "", err
	}
	acc, err := u.UserStorage.GetAccountByUsername(username)
	if err != nil {
		return "", ErrInvalidLogin
	}
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Credentials.Password), []byte(password)); err != nil {
		return "", ErrInvalidPassword
	}
	token, err := u.Auth.IssueToken(acc.Id)
	return Token(token), err
}

func (u DelegatedUserInterface) GetUser(user UserId) (User, error) {
	a, err := u.UserStorage.GetAccountById(uint(user))
	if err != nil {
		return User{}, err
	}

	return User{Id: UserId(a.Id), Username: a.Username}, nil
}

func (u DelegatedUserInterface) PostSnippet(author UserId, contents string, language ProgrammingLanguage) (Snippet, error) {
	id, err := u.SnippetStorage.AddSnippet(Snippet{
		Contents: contents,
		Language: language,
		Author:   author,
	})

	if err != nil {
		return Snippet{}, err
	}

	return u.SnippetStorage.GetSnippet(id)
}

func (u DelegatedUserInterface) GetSnippetsByUser(user UserId, current UserId) ([]Snippet, error) {
	s, err := u.SnippetStorage.GetSnippetsByUser(user)
	if err != nil {
		return []Snippet{}, err
	}

	for i := range s {
		vote, err := u.SnippetStorage.GetVote(current, s[i].Id)
		if err != nil {
			return []Snippet{}, err
		}
		s[i].CurrentUserVote = vote
	}

	return s, nil
}

func (u DelegatedUserInterface) GetSnippetsByLanguage(language ProgrammingLanguage, current UserId) ([]Snippet, error) {
	s, err := u.SnippetStorage.GetSnippetsByLanguage(language)
	if err != nil {
		return []Snippet{}, err
	}

	for i := range s {
		vote, err := u.SnippetStorage.GetVote(current, s[i].Id)
		if err != nil {
			return []Snippet{}, err
		}
		s[i].CurrentUserVote = vote
	}

	return s, nil
}

func (u DelegatedUserInterface) GetSnippet(current UserId, snippet SnippetId) (Snippet, error) {
	s, err := u.SnippetStorage.GetSnippet(snippet)
	if err != nil {
		return Snippet{}, err
	}

	vote, err := u.SnippetStorage.GetVote(current, s.Id)
	if err != nil {
		return Snippet{}, err
	}

	s.CurrentUserVote = vote
	return s, nil
}

func (u DelegatedUserInterface) DeleteSnippet(current UserId, snippet SnippetId) error {
	s, err := u.SnippetStorage.GetSnippet(snippet)
	if err != nil {
		return err
	}

	if s.Author != current {
		return MustBeSnippetAuthorErr
	}

	return u.SnippetStorage.DeleteSnippet(snippet)
}

func (u DelegatedUserInterface) Vote(current UserId, snippet SnippetId, vote int) error {
	return u.SnippetStorage.Vote(current, snippet, vote)
}

func (u DelegatedUserInterface) PostComment(author UserId, contents string, snippet SnippetId) (Comment, error) {
	id, err := u.SnippetStorage.AddComment(Comment{
		Contents: contents,
		Snippet:  snippet,
		Author:   author,
	})

	if err != nil {
		return Comment{}, err
	}

	return u.SnippetStorage.GetComment(id)
}

func (u DelegatedUserInterface) GetComments(snippet SnippetId) ([]Comment, error) {
	return u.SnippetStorage.GetComments(snippet)
}

func (u DelegatedUserInterface) DeleteComment(current UserId, comment CommentId) error {
	c, err := u.SnippetStorage.GetComment(comment)
	if err != nil {
		return err
	}

	if c.Author != current {
		return MustBeCommentAuthorErr
	}

	return u.SnippetStorage.DeleteComment(comment)
}
