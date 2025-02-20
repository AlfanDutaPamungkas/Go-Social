package cache

import (
	"context"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUsersStore{},
	}
}

type MockUsersStore struct{}

func (m *MockUsersStore) Get(context.Context, int64) (*store.User, error){
	return nil, nil
}

func (m *MockUsersStore) Set(context.Context, *store.User) error{
	return nil
}
