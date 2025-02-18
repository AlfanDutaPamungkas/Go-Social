package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type RoleStore struct {
	db *pgxpool.Pool
}

func (s *RoleStore) GetByName(ctx context.Context, slug string) (*Role, error) {
	query := `
		SELECT id, name, description, level
		FROM roles
		WHERE name = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var role Role
	err := s.db.QueryRow(ctx, query, slug).Scan(
		&role.ID, 
		&role.Name,
		&role.Description,
		&role.Level,
	)

	if err != nil {
		return nil, err
	}

	return &role, nil
}
