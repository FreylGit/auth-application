package postgresql

import (
	"auth-application/internal/domain/models"
	"auth-application/internal/storage"
	"context"
	"fmt"
	"time"
)

func (s *Storage) Refresh(ctx context.Context, userId int64, token string) (rToken models.RefreshToken, err error) {
	smtp, err := s.db.Prepare("SELECT id,user_id,token,expired_date FROM refresh_tokens WHERE user_id = $1 AND token = $2")
	if err != nil {
		return models.RefreshToken{}, fmt.Errorf("internal: %w", storage.ErrorSqlSyntax)
	}

	row := smtp.QueryRowContext(ctx, userId, token)
	var findToken string
	var id int64
	var user_id int64
	var expDate time.Time
	err = row.Scan(&id, &user_id, &findToken, &expDate)
	if err != nil {
		return models.RefreshToken{}, fmt.Errorf("internal: %w", storage.ErrorScan)
	}

	return models.RefreshToken{
		Id:      id,
		UserId:  user_id,
		Token:   findToken,
		ExpDate: expDate,
	}, nil
}

func (s *Storage) SaveRefresh(ctx context.Context, userId int64, token string) error {
	smtp, err := s.db.Prepare("INSERT INTO refresh_tokens (user_id,token,expired_date) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("internal: %w", storage.ErrorSqlSyntax)
	}
	expiredDate := time.Now().Add(time.Hour * 48)
	res, err := smtp.ExecContext(ctx, userId, token, expiredDate.Format("2006-01-02"))
	_ = res
	if err != nil {
		return fmt.Errorf("argument: %w", storage.ErrorUpdate)
	}

	return nil
}

func (s *Storage) UpdateRefresh(ctx context.Context, userId int64, newToken string, prevToken string) error {
	smtp, err := s.db.Prepare("UPDATE refresh_tokens SET token = $1 WHERE token = $2 AND user_id = $3")
	if err != nil {
		return fmt.Errorf("internal: %w", storage.ErrorSqlSyntax)
	}
	res, err := smtp.ExecContext(ctx, newToken, prevToken, userId)
	_ = res
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("argument: %w", storage.ErrorUpdate)
	}

	return nil
}
