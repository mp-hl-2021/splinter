package usecases

import (
	"errors"
	"github.com/mp-hl-2021/splinter/auth"
	"github.com/mp-hl-2021/splinter/highlighter"
	"github.com/mp-hl-2021/splinter/types"
	"golang.org/x/crypto/bcrypt"
)

var (
	MustBeSnippetAuthorErr = errors.New("must be snippet's author")
	MustBeCommentAuthorErr = errors.New("must be comment's author")
	InvalidVoteErr         = errors.New("invalid vote")
)

type DelegatedUserInterface struct {
	Auth           auth.Authenticator
	UserStorage    auth.UserStorage
	SnippetStorage types.SnippetStorage
	Highlighter    highlighter.Highlighter
}

func (u *DelegatedUserInterface) CreateAccount(username, password string) (types.User, error) {
	if err := auth.ValidateUsername(username); err != nil {
		return types.User{}, err
	}
	if err := auth.ValidatePassword(password); err != nil {
		return types.User{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return types.User{}, err
	}
	acc, err := u.UserStorage.CreateAccount(auth.Credentials{
		Username: username,
		Password: string(hashedPassword),
	})
	if err != nil {
		return types.User{}, err
	}
	return types.User{Id: types.UserId(acc.Id), Username: username}, nil
}

func (u *DelegatedUserInterface) Authenticate(username, password string) (types.Token, error) {
	if err := auth.ValidateUsername(username); err != nil {
		return "", err
	}
	if err := auth.ValidatePassword(password); err != nil {
		return "", err
	}
	acc, err := u.UserStorage.GetAccountByUsername(username)
	if err != nil {
		return "", types.ErrInvalidLogin
	}
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Credentials.Password), []byte(password)); err != nil {
		return "", types.ErrInvalidPassword
	}
	token, err := u.Auth.IssueToken(acc.Id)
	return types.Token(token), err
}

func (u DelegatedUserInterface) GetUser(user types.UserId) (types.User, error) {
	a, err := u.UserStorage.GetAccountById(uint(user))
	if err != nil {
		return types.User{}, err
	}

	return types.User{Id: types.UserId(a.Id), Username: a.Username}, nil
}

func (u DelegatedUserInterface) PostSnippet(author types.UserId, contents string, language types.ProgrammingLanguage) (types.Snippet, error) {
	id, err := u.SnippetStorage.AddSnippet(types.Snippet{
		Contents: contents,
		Language: language,
		Author:   author,
	})

	if err != nil {
		return types.Snippet{}, err
	}

	snippet, err := u.SnippetStorage.GetSnippet(id)

	if err != nil {
		return types.Snippet{}, err
	}

	u.Highlighter.Post(snippet)

	return snippet, nil
}

func (u DelegatedUserInterface) GetSnippetsByUser(user types.UserId, current types.UserId) ([]types.Snippet, error) {
	s, err := u.SnippetStorage.GetSnippetsByUser(user)
	if err != nil {
		return []types.Snippet{}, err
	}

	for i := range s {
		vote, err := u.SnippetStorage.GetVote(current, s[i].Id)
		if err != nil {
			return []types.Snippet{}, err
		}
		s[i].CurrentUserVote = vote
	}

	return s, nil
}

func (u DelegatedUserInterface) GetSnippetsByLanguage(language types.ProgrammingLanguage, current types.UserId) ([]types.Snippet, error) {
	s, err := u.SnippetStorage.GetSnippetsByLanguage(language)
	if err != nil {
		return []types.Snippet{}, err
	}

	for i := range s {
		vote, err := u.SnippetStorage.GetVote(current, s[i].Id)
		if err != nil {
			return []types.Snippet{}, err
		}
		s[i].CurrentUserVote = vote
	}

	return s, nil
}

func (u DelegatedUserInterface) GetSnippet(current types.UserId, snippet types.SnippetId) (types.Snippet, error) {
	s, err := u.SnippetStorage.GetSnippet(snippet)
	if err != nil {
		return types.Snippet{}, err
	}

	vote, err := u.SnippetStorage.GetVote(current, s.Id)
	if err != nil {
		return types.Snippet{}, err
	}

	s.CurrentUserVote = vote
	return s, nil
}

func (u DelegatedUserInterface) DeleteSnippet(current types.UserId, snippet types.SnippetId) error {
	s, err := u.SnippetStorage.GetSnippet(snippet)
	if err != nil {
		return err
	}

	if s.Author != current {
		return MustBeSnippetAuthorErr
	}

	return u.SnippetStorage.DeleteSnippet(snippet)
}

func (u DelegatedUserInterface) Vote(current types.UserId, snippet types.SnippetId, vote int) error {
	if vote < -1 || vote > 1 {
		return InvalidVoteErr
	}

	if _, err := u.SnippetStorage.GetSnippet(snippet); err != nil {
		return err
	}

	return u.SnippetStorage.Vote(current, snippet, vote)
}

func (u DelegatedUserInterface) PostComment(author types.UserId, contents string, snippet types.SnippetId) (types.Comment, error) {
	id, err := u.SnippetStorage.AddComment(types.Comment{
		Contents: contents,
		Snippet:  snippet,
		Author:   author,
	})

	if err != nil {
		return types.Comment{}, err
	}

	return u.SnippetStorage.GetComment(id)
}

func (u DelegatedUserInterface) GetComments(snippet types.SnippetId) ([]types.Comment, error) {
	return u.SnippetStorage.GetComments(snippet)
}

func (u DelegatedUserInterface) DeleteComment(current types.UserId, comment types.CommentId) error {
	c, err := u.SnippetStorage.GetComment(comment)
	if err != nil {
		return err
	}

	if c.Author != current {
		return MustBeCommentAuthorErr
	}

	return u.SnippetStorage.DeleteComment(comment)
}
