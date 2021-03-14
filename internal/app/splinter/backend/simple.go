package backend

import (
	"fmt"
	"github.com/mp-hl-2021/splinter/internal/app/splinter/model"
)

type SimpleUserInterface struct {
	users    []model.User
	snippets []model.Snippet
	lastId   int64
}

func NewSimpleUserInterface() *SimpleUserInterface {
	return &SimpleUserInterface{
		users:    make([]model.User, 0),
		snippets: make([]model.Snippet, 0),
		lastId:   0,
	}
}

func (x *SimpleUserInterface) getUser(userId string) *model.User {
	for _, u := range x.users {
		if u.Id == userId {
			return &u
		}
	}
	return nil
}

func (x *SimpleUserInterface) getSnippet(snippetId string) *model.Snippet {
	for _, s := range x.snippets {
		if s.Id == snippetId {
			return &s
		}
	}
	return nil
}

func (x *SimpleUserInterface) nextId() string {
	x.lastId += 1
	return fmt.Sprintf("%d", x.lastId)
}

func (x *SimpleUserInterface) Register(request model.RegisterRequest) (*model.AuthenticateResponse, error) {
	for _, u := range x.users {
		if u.Username == request.Username {
			return nil, fmt.Errorf("user @%s already exists", request.Username)
		}
	}

	user := model.User{
		Id:       x.nextId(),
		Username: request.Username,
	}

	x.users = append(x.users, user)
	return &model.AuthenticateResponse{
		Token: model.Token(user.Id),
		User:  user,
	}, nil
}

func (x *SimpleUserInterface) Authenticate(request model.AuthenticateRequest) (*model.AuthenticateResponse, error) {
	for _, u := range x.users {
		if u.Username == request.Username {
			return &model.AuthenticateResponse{
				Token: model.Token(u.Id),
				User:  u,
			}, nil
		}
	}
	return nil, fmt.Errorf("invalid username")
}

func (x *SimpleUserInterface) GetUserById(token model.Token, id string) (*model.User, error) {
	u := x.getUser(id)
	if u != nil {
		return u, nil
	}
	return nil, fmt.Errorf("user doesn't exist")
}

func (x *SimpleUserInterface) GetSnippetById(token model.Token, id string) (*model.Snippet, error) {
	s := x.getSnippet(id)
	if s != nil {
		return s, nil
	}
	return nil, fmt.Errorf("snippet doesn't exist")
}

func (x *SimpleUserInterface) GetSnippetsByUser(token model.Token, userId string) ([]model.Snippet, error) {
	var snippets []model.Snippet
	for _, s := range x.snippets {
		if s.AuthorId == userId {
			snippets = append(snippets, s)
		}
	}
	return snippets, nil
}

func (x *SimpleUserInterface) GetSnippetsByLanguage(token model.Token, language model.ProgrammingLanguage) ([]model.Snippet, error) {
	var snippets []model.Snippet
	for _, s := range x.snippets {
		if s.Language == language {
			snippets = append(snippets, s)
		}
	}
	return snippets, nil
}

func (x *SimpleUserInterface) GetSnippetsFeed(token model.Token) ([]model.Snippet, error) {
	return x.snippets, nil
}

func (x *SimpleUserInterface) PostSnippet(token model.Token, request model.PostSnippetRequest) (*model.Snippet, error) {
	s := model.Snippet{
		Id:       x.nextId(),
		Contents: request.Contents,
		Language: request.Language,
		AuthorId: string(token),
		Rating:   model.Rating{},
	}
	x.snippets = append(x.snippets, s)
	return &s, nil
}

func (x *SimpleUserInterface) DeleteSnippetById(token model.Token, id string) error {
	return fmt.Errorf("not implemented")
}

func (x *SimpleUserInterface) VoteSnippet(token model.Token, request model.VoteSnippetRequest) error {
	return fmt.Errorf("not implemented")
}

func (x *SimpleUserInterface) PostComment(token model.Token, snippetId string, request model.PostCommentRequest) (*model.Comment, error) {
	return nil, fmt.Errorf("not implemented")
}

func (x *SimpleUserInterface) GetCommentsBySnippetId(token model.Token, snippetId string) ([]model.Comment, error) {
	return nil, fmt.Errorf("not implemented")
}

func (x *SimpleUserInterface) GetSubscriptions(token model.Token) ([]model.Subscription, error) {
	return nil, fmt.Errorf("not implemented")
}

func (x *SimpleUserInterface) SubscribeToLanguage(token model.Token, language model.ProgrammingLanguage) error {
	return fmt.Errorf("not implemented")
}

func (x *SimpleUserInterface) UnsubscribeFromLanguage(token model.Token, language model.ProgrammingLanguage) error {
	return fmt.Errorf("not implemented")
}

func (x *SimpleUserInterface) SubscribeToUser(token model.Token, userId string) error {
	return fmt.Errorf("not implemented")
}

func (x *SimpleUserInterface) UnsubscribeFromUser(token model.Token, userId string) error {
	return fmt.Errorf("not implemented")
}
