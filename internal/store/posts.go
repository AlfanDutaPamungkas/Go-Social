package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comments_count"`
}

type PostStore struct {
	db *pgxpool.Pool
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		post.Tags,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.Tags,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, nil
		}
	}

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, tags = $3, updated_at = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		post.Tags,
		post.UpdatedAt,
		post.ID,
		post.Version,
	).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, p *PaginatedFeedQuery) ([]PostWithMetadata, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1 // Mulai dari $1 untuk parameter pertama

	// Query dasar
	query := `
	SELECT 
		p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
		u.username,
		count(c.id) AS comments_count
	FROM posts p
	LEFT JOIN comments c ON c.post_id = p.id
	LEFT JOIN users u ON p.user_id = u.id
	LEFT JOIN followers f ON f.follower_id = p.user_id
	WHERE (f.user_id = $1 OR p.user_id = $1)
	`

	args = append(args, userID)

	// Filter Search
	if p.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"to_tsvector('english', p.title || ' ' || p.content) @@ plainto_tsquery('english', $%d::text)",
			argIndex+1))
		args = append(args, p.Search)
		argIndex++
	}	

	// Filter Tags
	if len(p.Tags) > 0 {
		conditions = append(conditions, fmt.Sprintf("p.tags @> $%d", argIndex+1))
		args = append(args, p.Tags)
		argIndex++
	}

	// Filter Since
	if !p.Since.IsZero() {
		conditions = append(conditions, fmt.Sprintf("p.created_at >= $%d", argIndex+1))
		args = append(args, p.Since)
		argIndex++
	}

	// Filter Until
	if !p.Until.IsZero() {
		conditions = append(conditions, fmt.Sprintf("p.created_at <= $%d", argIndex+1))
		args = append(args, p.Until)
		argIndex++
	}

	// Tambahkan kondisi ke query jika ada filter
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Tambahkan ORDER, LIMIT, dan OFFSET
	query += fmt.Sprintf(" GROUP BY p.id, u.username ORDER BY p.created_at %s LIMIT $%d OFFSET $%d", p.Sort, argIndex+1, argIndex+2)
	args = append(args, p.Limit, p.Offset)

	// Eksekusi Query
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan hasil ke dalam struct
	var feeds []PostWithMetadata
	for rows.Next() {
		var feed PostWithMetadata
		if err := rows.Scan(
			&feed.ID,
			&feed.UserID,
			&feed.Title,
			&feed.Content,
			&feed.CreatedAt,
			&feed.Version,
			&feed.Tags,
			&feed.User.Username,
			&feed.CommentCount,
		); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	fmt.Println("FINAL QUERY:", query)
	fmt.Println("ARGS:", args)	

	return feeds, nil
}
