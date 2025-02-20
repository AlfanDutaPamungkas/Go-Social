package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Create(ctx context.Context, tx pgx.Tx, u *User) error {
	return nil
}

func (m *MockUserStore) GetByID(context.Context, int64) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) GetByEmail(context.Context, string) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) CreateAndInvite(context.Context, *User, string, time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(context.Context, string) error {
	return nil
}

func (m *MockUserStore) Delete(context.Context, int64) error {
	return nil
}
