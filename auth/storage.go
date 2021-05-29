package auth

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
)

type Account struct {
	Id uint
	Credentials
}

type Credentials struct {
	Username string
	Password string
}

type UserStorage interface {
	CreateAccount(cred Credentials) (Account, error)
	GetAccountById(id uint) (Account, error)
	GetAccountByUsername(username string) (Account, error)
}
