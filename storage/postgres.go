package storage

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mp-hl-2021/splinter/auth"
	"github.com/mp-hl-2021/splinter/usecases"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(connStr string) (Postgres, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return Postgres{}, err
	}

	return Postgres{db}, nil
}

func (p Postgres) Close() error {
	return p.db.Close()
}

func (p Postgres) AddSnippet(snippet usecases.Snippet) (usecases.SnippetId, error) {
	res, err := p.db.Exec(`
	insert into snippet (contents, language, author, createdAt) values
	($1, $2, $3, now());
`, snippet.Contents, snippet.Language, snippet.Author)
	if err != nil {
		return usecases.SnippetId(0), err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return usecases.SnippetId(0), err
	}

	return usecases.SnippetId(id), nil
}

func (p Postgres) GetSnippetsByUser(user usecases.UserId) ([]usecases.Snippet, error) {
	rows, err := p.db.Query(`
select id, contents, language, author, likes, dislikes, createdAt from snippet where author = $1
`, user)

	if err != nil {
		return []usecases.Snippet{}, err
	}

	var res []usecases.Snippet

	for rows.Next() {
		var s usecases.Snippet
		err = rows.Scan(&s.Id, &s.Contents, &s.Language, &s.Rating.Likes, &s.Rating.Dislikes, &s.CreatedAt)
		if err != nil {
			return []usecases.Snippet{}, err
		}
		res = append(res, s)
	}

	return res, nil
}

func (p Postgres) GetSnippetsByLanguage(language usecases.ProgrammingLanguage) ([]usecases.Snippet, error) {
	rows, err := p.db.Query(`
select id, contents, language, author, likes, dislikes, createdAt from snippet where language = $1
`, language)

	if err != nil {
		return []usecases.Snippet{}, err
	}

	var res []usecases.Snippet

	for rows.Next() {
		var s usecases.Snippet
		err = rows.Scan(&s.Id, &s.Contents, &s.Language, &s.Rating.Likes, &s.Rating.Dislikes, &s.CreatedAt)
		if err != nil {
			return []usecases.Snippet{}, err
		}
		res = append(res, s)
	}

	return res, nil
}

func (p Postgres) GetSnippet(snippet usecases.SnippetId) (usecases.Snippet, error) {
	row := p.db.QueryRow(`
select id, contents, language, author, likes, dislikes, createdAt from snippet where id = $1
`, snippet)

	var s usecases.Snippet
	err := row.Scan(&s.Id, &s.Contents, &s.Language, &s.Rating.Likes, &s.Rating.Dislikes, &s.CreatedAt)
	if err != nil {
		return usecases.Snippet{}, err
	}

	return s, nil
}

func (p Postgres) DeleteSnippet(snippet usecases.SnippetId) error {
	_, err := p.db.Exec(`
delete from snippet where id = $1
`, snippet)
	return err
}

func (p Postgres) Vote(user usecases.UserId, snippet usecases.SnippetId, vote int) error {
	_, err := p.db.Exec(`
insert into vote (snippet, "user", vote) values ($1, $2, $3)
on conflict on constraint vote_pkey do update
set vote = $3;
`, snippet, user, vote)
	return err
}

func (p Postgres) GetVote(user usecases.UserId, snippet usecases.SnippetId) (int, error) {
	row := p.db.QueryRow(`
select vote from vote where snippet = $1 and user = $2;
`, snippet, user)
	var vote int
	err := row.Scan(&vote)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return vote, nil
}

func (p Postgres) AddComment(comment usecases.Comment) (usecases.CommentId, error) {
	res, err := p.db.Exec(`
insert into comment (contents, snippet, author, createdAt) values ($1, $2, $3, now());
`, comment.Contents, comment.Snippet, comment.Author)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return usecases.CommentId(id), nil
}

func (p Postgres) GetComment(comment usecases.CommentId) (usecases.Comment, error) {
	row := p.db.QueryRow(`
select id, contents, snippet, author, createdAt from comment where id = $1
`, comment)

	var c usecases.Comment
	err := row.Scan(&c.Id, &c.Contents, &c.Snippet, &c.Author, &c.CreatedAt)
	if err != nil {
		return usecases.Comment{}, err
	}

	return c, nil
}

func (p Postgres) GetComments(snippet usecases.SnippetId) ([]usecases.Comment, error) {
	rows, err := p.db.Query(`
select id, contents, snippet, author, createdAt from comment where snippet = $1
`, snippet)

	if err != nil {
		return []usecases.Comment{}, err
	}

	var res []usecases.Comment

	for rows.Next() {
		var c usecases.Comment
		err = rows.Scan(&c.Id, &c.Contents, &c.Snippet, &c.Author, &c.CreatedAt)
		if err != nil {
			return []usecases.Comment{}, err
		}
		res = append(res, c)
	}

	return res, nil
}

func (p Postgres) DeleteComment(comment usecases.CommentId) error {
	_, err := p.db.Exec(`
delete from comment where id = $1
`, comment)
	return err
}

func (p Postgres) CreateAccount(cred auth.Credentials) (auth.Account, error) {
	res, err := p.db.Exec(`
insert into "user" (username, password) values ($1, $2)
`, cred.Username, cred.Password)
	if err != nil {
		return auth.Account{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return auth.Account{}, err
	}

	return auth.Account{
		Id:          uint(id),
		Credentials: cred,
	}, nil
}

func (p Postgres) GetAccountById(id uint) (auth.Account, error) {
	row := p.db.QueryRow(`
select id, username, password from "user" where id = $1
`, id)
	var a auth.Account
	err := row.Scan(&a.Id, &a.Username, &a.Password)
	if err != nil {
		return auth.Account{}, err
	}
	return a, nil
}

func (p Postgres) GetAccountByUsername(username string) (auth.Account, error) {
	row := p.db.QueryRow(`
select id, username, password from "user" where username = $1
`, username)
	var a auth.Account
	err := row.Scan(&a.Id, &a.Username, &a.Password)
	if err != nil {
		return auth.Account{}, err
	}
	return a, nil
}
