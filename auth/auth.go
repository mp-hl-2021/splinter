package auth

type Authenticator interface {
	IssueToken(userId uint) (string, error)
	UserIdByToken(token string) (uint, error)
}