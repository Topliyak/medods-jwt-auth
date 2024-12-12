package auth

import (
	"context"

	"github.com/jackc/pgx/v5"

	passwordUtils "github.com/medods-jwt-auth/utils/password"
	userCrud "github.com/medods-jwt-auth/utils/crud/user"
)

func CheckCredentials(ctx context.Context, tx pgx.Tx, email, password string) bool {
	user := userCrud.GetUserByEmail(ctx, tx, email)
	if user == nil {
		return false
	}
	return passwordUtils.ValidatePassword(password, user.Password)
}
