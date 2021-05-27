package storage

import (
	"errors"
	"github.com/mp-hl-2021/splinter/auth"
	"github.com/mp-hl-2021/splinter/usecases"
	"sync"
	"time"
)

var (
	NoSuchSnippetErr = errors.New("no such snippet")
	NoSuchCommentErr = errors.New("no such comment")
	InvalidVoteErr   = errors.New("invalid vote")
)

type SnippetVote struct {
	UserId    usecases.UserId
	SnippetId usecases.SnippetId
	Vote      int
}

type Memory struct {
	snippets           []usecases.Snippet
	votes              []SnippetVote
	comments           []usecases.Comment
	accountsById       map[uint]auth.Account
	accountsByUsername map[string]auth.Account
	nextId             uint
	mu                 *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		accountsById:       make(map[uint]auth.Account),
		accountsByUsername: make(map[string]auth.Account),
		mu:                 &sync.Mutex{},
	}
}

func (m *Memory) CreateAccount(cred auth.Credentials) (auth.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.accountsByUsername[cred.Username]; ok {
		return auth.Account{}, auth.ErrAlreadyExist
	}
	a := auth.Account{
		Id:          m.nextId,
		Credentials: cred,
	}
	m.accountsById[a.Id] = a
	m.accountsByUsername[a.Username] = a
	m.nextId++
	return a, nil
}

func (m *Memory) AddSnippet(snippet usecases.Snippet) (usecases.SnippetId, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	snippet.Id = usecases.SnippetId(m.nextId)
	snippet.CreatedAt = time.Now()
	m.nextId++
	m.snippets = append(m.snippets, snippet)
	return snippet.Id, nil
}

func (m *Memory) GetAccountById(id uint) (auth.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsById[id]
	if !ok {
		return a, auth.ErrNotFound
	}
	return a, nil
}

func (m *Memory) GetAccountByUsername(username string) (auth.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsByUsername[username]
	if !ok {
		return a, auth.ErrNotFound
	}
	return a, nil
}

func (m *Memory) GetSnippetsByUser(user usecases.UserId) ([]usecases.Snippet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var res []usecases.Snippet
	for _, s := range m.snippets {
		if s.Author == user {
			res = append(res, s)
		}
	}
	return res, nil
}

func (m *Memory) GetSnippetsByLanguage(language usecases.ProgrammingLanguage) ([]usecases.Snippet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var res []usecases.Snippet
	for _, s := range m.snippets {
		if s.Language == language {
			res = append(res, s)
		}
	}
	return res, nil
}

func (m *Memory) GetSnippet(snippet usecases.SnippetId) (usecases.Snippet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.snippets {
		if s.Id == snippet {
			return s, nil
		}
	}
	return usecases.Snippet{}, NoSuchSnippetErr
}

func (m *Memory) DeleteSnippet(snippet usecases.SnippetId) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var res []usecases.Snippet
	for _, s := range m.snippets {
		if s.Id != snippet {
			res = append(res, s)
		}
	}
	m.snippets = res
	return nil
}

func (m *Memory) Vote(user usecases.UserId, snippet usecases.SnippetId, vote int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if vote != -1 && vote != 0 && vote != 1 {
		return InvalidVoteErr
	}

	for i, v := range m.votes {
		if v.UserId == user && v.SnippetId == snippet {
			m.votes[i].Vote = vote
			return nil
		}
	}

	m.votes = append(m.votes, SnippetVote{
		UserId:    user,
		SnippetId: snippet,
		Vote:      vote,
	})

	return nil
}

func (m *Memory) GetVote(user usecases.UserId, snippet usecases.SnippetId) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, v := range m.votes {
		if v.UserId == user && v.SnippetId == snippet {
			return v.Vote, nil
		}
	}
	return 0, nil
}

func (m *Memory) AddComment(comment usecases.Comment) (usecases.CommentId, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	comment.Id = usecases.CommentId(m.nextId)
	comment.CreatedAt = time.Now()
	m.nextId++
	m.comments = append(m.comments, comment)
	return comment.Id, nil
}

func (m *Memory) GetComment(comment usecases.CommentId) (usecases.Comment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.comments {
		if c.Id == comment {
			return c, nil
		}
	}

	return usecases.Comment{}, NoSuchCommentErr
}

func (m *Memory) GetComments(snippet usecases.SnippetId) ([]usecases.Comment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var res []usecases.Comment
	for _, c := range m.comments {
		if c.Snippet == snippet {
			res = append(res, c)
		}
	}

	return res, nil
}

func (m *Memory) DeleteComment(comment usecases.CommentId) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var res []usecases.Comment
	for _, c := range m.comments {
		if c.Id != comment {
			res = append(res, c)
		}
	}
	m.comments = res
	return nil
}
