package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
	}

	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *pgxpool.Pool) Storage{
	return Storage{
		Posts: &PostStore{db},
		Users: &UsersStore{db},
	}
}
