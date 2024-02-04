package postgresql

import (
	"auth-application/internal/domain/models"
	"auth-application/internal/storage"
	"context"
	"fmt"
	"github.com/lib/pq"
)

func (s *Storage) SaveUser(ctx context.Context, newUser models.User) (userId int64, err error) {
	query := `
        INSERT INTO users(email, password_hash, name)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	var id int64
	err = s.db.QueryRowContext(ctx, query, newUser.Email, newUser.PasswordHash, newUser.Name).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, fmt.Errorf("argument email already exists %w", storage.ErrorNoUnique)
			}
		}
		return 0, fmt.Errorf("error inserting user: %w", storage.ErrorSave)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	smtp, err := s.db.Prepare("SELECT id, email, password_hash, name,  create_date " +
		"FROM users WHERE email = $1")
	if err != nil {
		return nil, fmt.Errorf("internal: %w", storage.ErrorSqlSyntax)
	}

	row := smtp.QueryRowContext(ctx, email)
	var user models.User
	err = row.Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Name, &user.CreateDate)
	if err != nil {
		return nil, fmt.Errorf("argument: %w", storage.ErrNotFound)
	}

	return &user, nil
}

func (s *Storage) UserById(ctx context.Context, id int64) (*models.User, error) {
	smtp, err := s.db.Prepare("SELECT id, email, password_hash, name, create_date FROM users WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("internal: %w", storage.ErrorSqlSyntax)
	}
	row := smtp.QueryRowContext(ctx, id)
	var user models.User
	err = row.Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Name, &user.CreateDate)
	if err != nil {
		return nil, fmt.Errorf("argument: %w", storage.ErrNotFound)
	}

	return &user, nil
}
