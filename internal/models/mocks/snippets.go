package mocks

import (
	"github.com/nightfaust/snippetbox/internal/models"
	"time"
)

var mockSnippet = models.Snippet{
	Created: time.Now(),
	Expires: time.Now(),
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(string, string, int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}
