package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
	}

	Users interface {
		Create(context.Context, *User) error
	}

	Comments interface {
		GetCommentsByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
}

func NewStorage(db *pgxpool.Pool) Storage{
	return Storage{
		Posts: &PostStore{db},
		Users: &UsersStore{db},
		Comments: &CommentStore{db},
	}
}
