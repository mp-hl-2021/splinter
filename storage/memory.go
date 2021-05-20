package storage

import (
	"sync"
)

type Memory struct {
	accountsById       map[uint]Account
	accountsByUsername map[string]Account
	nextId             uint
	mu                 *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		accountsById:       make(map[uint]Account),
		accountsByUsername: make(map[string]Account),
		mu:                 &sync.Mutex{},
	}
}

func (m *Memory) CreateAccount(cred Credentials) (Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.accountsByUsername[cred.Username]; ok {
		return Account{}, ErrAlreadyExist
	}
	a := Account{
		Id:          m.nextId,
		Credentials: cred,
	}
	m.accountsById[a.Id] = a
	m.accountsByUsername[a.Username] = a
	m.nextId++
	return a, nil
}

func (m *Memory) GetAccountById(id uint) (Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsById[id]
	if !ok {
		return a, ErrNotFound
	}
	return a, nil
}

func (m *Memory) GetAccountByUsername(username string) (Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsByUsername[username]
	if !ok {
		return a, ErrNotFound
	}
	return a, nil
}
