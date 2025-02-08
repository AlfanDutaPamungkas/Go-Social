package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64, *PaginatedFeedQuery) ([]PostWithMetadata, error)
	}

	Users interface {
		GetByID(context.Context, int64) (*User, error)
		Create(context.Context, *User) error
	}

	Comments interface {
		Create(context.Context, *Comment) error
		GetCommentsByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}

	Followers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		Unfollow(ctx context.Context, followerID, userID int64) error
	}
}

func NewStorage(db *pgxpool.Pool) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
