package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
}

// Insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires)
  	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	var s Snippet
	query := `SELECT id, title, content, created, expires FROM snippets
  WHERE expires > UTC_TIMESTAMP() AND id = ?`

	err := m.DB.QueryRow(query, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// Latest return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	var snippets []Snippet
	query := `SELECT id, title, content, created, expires FROM snippets
  WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

type Snippet struct {
	Created time.Time
	Expires time.Time
	ID      int
	Title   string
	Content string
}
type SnippetModel struct {
	DB *sql.DB
}
