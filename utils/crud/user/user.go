package user

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/medods-jwt-auth/models"
)

func GetUserById(ctx context.Context, tx pgx.Tx, id int64) *models.User {
	sql := "SELECT id, email, password, refresh, refresh_issued_at FROM users WHERE id = $1"
	row := tx.QueryRow(ctx, sql, id)
	user, err := getUserFromRow(row)
	if err != nil {
		return nil
	}
	return user
}

func GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) *models.User {
	sql := "SELECT id, email, password, refresh, refresh_issued_at FROM users WHERE email = $1"
	row := tx.QueryRow(ctx, sql, email)
	user, err := getUserFromRow(row)
	if err == pgx.ErrNoRows {
		return nil
	}
	return user
}

func UpdateRefreshTokenByUserId(ctx context.Context, tx pgx.Tx, id int64, refresh string, issuedAt time.Time) (err error) {
	sql := `
		UPDATE users
		SET
			refresh = $1,
			refresh_issued_at = $2
		WHERE id = $3
	`
	_, err = tx.Exec(ctx, sql, refresh, issuedAt, id)
	return err
}

func UpdateRefreshTokenByUserEmail(ctx context.Context, tx pgx.Tx, email, refresh string, issuedAt time.Time) (err error) {
	sql := `
		UPDATE users
		SET
			refresh = $1,
			refresh_issued_at = $2
		WHERE email = $3
	`
	_, err = tx.Exec(ctx, sql, refresh, issuedAt, email)
	return err
}

func CreateUser(ctx context.Context, tx pgx.Tx, email, password string) error {
	sql := "INSERT INTO users (email, password) VALUES ($1, $2)"
	_, err := tx.Exec(ctx, sql, email, password)
	return err
}

func DeleteUserById(ctx context.Context, tx pgx.Tx, id int64) error {
	sql := "DELETE FROM users WHERE id = $1"
	_, err := tx.Exec(ctx, sql, id)
	return err
}

func getUserFromRow(row pgx.Row) (*models.User, error) {
	var user models.User
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Refresh, &user.RefreshIssuedAt)
	return &user, err
}
