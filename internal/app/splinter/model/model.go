package model

type ProgrammingLanguage string

type User struct {
	Id       string // Unique identifier, persists through username changes
	Username string // Visible username, can be changed
}

type SubscriptionType int

const (
	SubscriptionToLanguage SubscriptionType = iota
	SubscriptionToUser
)

type Subscription interface {
	Type() SubscriptionType
	Language() *ProgrammingLanguage
	User() *User
}

type Vote int

const (
	VoteUp   Vote = +1
	VoteDown Vote = -1
)

type Rating struct {
	VotesUp   int
	VotesDown int
}

type Snippet struct {
	Id       string
	Contents string
	Language ProgrammingLanguage
	AuthorId string
	Rating   Rating
}

type Comment struct {
	SnippetId string
	Contents  string
	AuthorId  string
}
