package storage

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mp-hl-2021/splinter/auth"
	"github.com/mp-hl-2021/splinter/types"
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

func (p Postgres) AddSnippet(snippet types.Snippet) (types.SnippetId, error) {
	var id int
	err := p.db.QueryRow(`
	insert into snippet (contents, language, author, createdAt) values
	($1, $2, $3, now()) returning (id);
`, snippet.Contents, snippet.Language, snippet.Author).Scan(&id)
	if err != nil {
		return types.SnippetId(0), err
	}

	return types.SnippetId(id), nil
}

func (p Postgres) SetSnippetHighlight(snippet types.SnippetId, highlight string) error {
	_, err := p.db.Exec(`
	update snippet set highlighted = $1 where id = $2;
`, highlight, snippet)
	return err
}

func (p Postgres) GetSnippetsByUser(user types.UserId) ([]types.Snippet, error) {
	rows, err := p.db.Query(`
select id, contents, highlighted, language, author, likes, dislikes, createdAt from snippet where author = $1
`, user)

	if err != nil {
		return []types.Snippet{}, err
	}

	var res []types.Snippet

	for rows.Next() {
		var s types.Snippet
		err = rows.Scan(&s.Id, &s.Contents, &s.HighlightedContents, &s.Language, &s.Author, &s.Rating.Likes, &s.Rating.Dislikes, &s.CreatedAt)
		if err != nil {
			return []types.Snippet{}, err
		}
		res = append(res, s)
	}

	return res, nil
}

func (p Postgres) GetSnippetsByLanguage(language types.ProgrammingLanguage) ([]types.Snippet, error) {
	rows, err := p.db.Query(`
select id, contents, highlighted, language, author, likes, dislikes, createdAt from snippet where language = $1
`, language)

	if err != nil {
		return []types.Snippet{}, err
	}

	var res []types.Snippet

	for rows.Next() {
		var s types.Snippet
		err = rows.Scan(&s.Id, &s.Contents, &s.HighlightedContents, &s.Language, &s.Author, &s.Rating.Likes, &s.Rating.Dislikes, &s.CreatedAt)
		if err != nil {
			return []types.Snippet{}, err
		}
		res = append(res, s)
	}

	return res, nil
}

func (p Postgres) GetSnippet(snippet types.SnippetId) (types.Snippet, error) {
	row := p.db.QueryRow(`
select id, contents, highlighted, language, author, likes, dislikes, createdAt from snippet where id = $1
`, snippet)

	var s types.Snippet
	err := row.Scan(&s.Id, &s.Contents, &s.HighlightedContents, &s.Language, &s.Author, &s.Rating.Likes, &s.Rating.Dislikes, &s.CreatedAt)
	if err != nil {
		return types.Snippet{}, err
	}

	return s, nil
}

func (p Postgres) DeleteSnippet(snippet types.SnippetId) error {
	_, err := p.db.Exec(`
delete from snippet where id = $1
`, snippet)
	return err
}

func (p Postgres) Vote(user types.UserId, snippet types.SnippetId, vote int) error {
	_, err := p.db.Exec(`
insert into vote (snippet, "user", vote) values ($1, $2, $3)
on conflict on constraint vote_pkey do update
set vote = $3;
`, snippet, user, vote)
	return err
}

func (p Postgres) GetVote(user types.UserId, snippet types.SnippetId) (int, error) {
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

func (p Postgres) AddComment(comment types.Comment) (types.CommentId, error) {
	var id int
	err := p.db.QueryRow(`
insert into comment (contents, snippet, author, createdAt) values ($1, $2, $3, now()) returning (id);
`, comment.Contents, comment.Snippet, comment.Author).Scan(&id)
	if err != nil {
		return 0, err
	}

	return types.CommentId(id), nil
}

func (p Postgres) GetComment(comment types.CommentId) (types.Comment, error) {
	row := p.db.QueryRow(`
select id, contents, snippet, author, createdAt from comment where id = $1
`, comment)

	var c types.Comment
	err := row.Scan(&c.Id, &c.Contents, &c.Snippet, &c.Author, &c.CreatedAt)
	if err != nil {
		return types.Comment{}, err
	}

	return c, nil
}

func (p Postgres) GetComments(snippet types.SnippetId) ([]types.Comment, error) {
	rows, err := p.db.Query(`
select id, contents, snippet, author, createdAt from comment where snippet = $1
`, snippet)

	if err != nil {
		return []types.Comment{}, err
	}

	var res []types.Comment

	for rows.Next() {
		var c types.Comment
		err = rows.Scan(&c.Id, &c.Contents, &c.Snippet, &c.Author, &c.CreatedAt)
		if err != nil {
			return []types.Comment{}, err
		}
		res = append(res, c)
	}

	return res, nil
}

func (p Postgres) DeleteComment(comment types.CommentId) error {
	_, err := p.db.Exec(`
delete from comment where id = $1
`, comment)
	return err
}

func (p Postgres) CreateAccount(cred auth.Credentials) (auth.Account, error) {
	var id int
	err := p.db.QueryRow(`
insert into "user" (username, password) values ($1, $2) returning (id)
`, cred.Username, cred.Password).Scan(&id)
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
	var a auth.Account
	err := p.db.QueryRow(`
select id, username, password from "user" where username = $1
`, username).Scan(&a.Id, &a.Username, &a.Password)
	if err != nil {
		return auth.Account{}, err
	}
	return a, nil
}
