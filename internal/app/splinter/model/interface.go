package model

type UserInterface interface {
	Register(request RegisterRequest) (*AuthenticateResponse, error)
	Authenticate(request AuthenticateRequest) (*AuthenticateResponse, error)

	GetSnippetById(token Token, id string) (*Snippet, error)
	GetSnippetsByUser(token Token, userId string) ([]Snippet, error)
	GetSnippetsByLanguage(token Token, language ProgrammingLanguage) ([]Snippet, error)
	GetSnippetsFeed(token Token) ([]Snippet, error)
	PostSnippet(token Token, request PostSnippetRequest) (*Snippet, error)
	DeleteSnippetById(token Token, id string) error
	VoteSnippet(token Token, request VoteSnippetRequest) error

	PostComment(token Token, request PostCommentRequest) (*Comment, error)
	GetCommentsBySnippetId(token Token, snippetId string) ([]Comment, error)

	GetSubscriptions(token Token) ([]Subscription, error)
	SubscribeToLanguage(token Token, language ProgrammingLanguage) error
	UnsubscribeFromLanguage(token Token, language ProgrammingLanguage) error
	SubscribeToUser(token Token, userId string) error
	UnsubscribeFromUser(token Token, userId string) error
}

type Token string

const UnauthenticatedToken Token = ""

type RegisterRequest struct {
	Username string
	Password string
}

type AuthenticateRequest struct {
	Username string
	Password string
}

type AuthenticateResponse struct {
	Token Token
	User  User
}

type PostSnippetRequest struct {
	Contents string
	Language ProgrammingLanguage
}

type VoteSnippetRequest struct {
	SnippetId string
	Vote      Vote
}

type PostCommentRequest struct {
	SnippetId string
	Contents  string
}
