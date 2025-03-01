package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Followers struct {
	UserID     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
	CreatedAt  int64 `json:"created_at"`
}

type FollowerStore struct {
	db *pgxpool.Pool
}

func (s *FollowerStore) Follow(ctx context.Context, userID, followerID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.Exec(ctx, query, userID, followerID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return ErrConflict
		}
	}

	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, userID, followerID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.Exec(ctx, query, userID, followerID)
	return err
}
