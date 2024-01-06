package mocks

import (
	"github.com/nightfaust/snippetbox/internal/models"
	"time"
)

var mockUser = models.User{
	Created:        time.Now(),
	ID:             1,
	Name:           "alice",
	Email:          "test@test.com",
	HashedPassword: nil,
}

type UserModel struct{}

func (m *UserModel) Insert(_ string, email string, _ string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Get(id int) (models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return models.User{}, models.ErrNoRecord
	}
}
