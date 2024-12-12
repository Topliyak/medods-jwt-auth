package token

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-randomstring/randomstring"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"

	"github.com/medods-jwt-auth/config"
	"github.com/medods-jwt-auth/db"
	userCrud "github.com/medods-jwt-auth/utils/crud/user"
	passwordUtils "github.com/medods-jwt-auth/utils/password"
	mailUtils "github.com/medods-jwt-auth/utils/mail"
)

const (
    REFRESH_TOKEN_LENGTH = 50
)

var (
    SECRET_KEY = []byte(config.SECRET_KEY)
)

type Claims struct {
	jwt.RegisteredClaims
    UserId int64 `json:"user_id"`
    IP string `json:"ip"`
}

func CreateToken(ctx context.Context, tx pgx.Tx, email string, ip string) (jwt_token, refresh string, err error)  {
    userId := userCrud.GetUserByEmail(ctx, tx, email).Id
    issuedAt := time.Unix(time.Now().Unix(), 0).UTC()

    if jwt_token, err = createJWT(userId, ip, issuedAt); err != nil {
        return "", "", err
    }

    if refresh, err = createRefresh(ctx, tx, userId, issuedAt); err != nil {
        return "", "", err
    }

    return jwt_token, refresh, nil
}

func createJWT(userId int64, ip string, issuedAt time.Time) (jwt_token string, err error) {
    expireAt := issuedAt.Add(config.JWT_DURATION)
    claims := &Claims{
        RegisteredClaims: jwt.RegisteredClaims{
            IssuedAt: jwt.NewNumericDate(issuedAt),
            ExpiresAt: jwt.NewNumericDate(expireAt),
        },
        UserId: userId,
        IP: ip,
    }

    jwt_token, err = jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(config.SECRET_KEY))
    if err != nil {
        return "", err
    }

    return jwt_token, nil
}

func createRefresh(ctx context.Context, tx pgx.Tx, userId int64, issuedAt time.Time) (refresh string, err error) {
    refresh = randomstring.Generate(
        REFRESH_TOKEN_LENGTH,
        randomstring.Digits,
        randomstring.Symbols,
        randomstring.LowerLetters,
        randomstring.UpperLetters,
    )

    var refreshHashed string
    if refreshHashed, err = passwordUtils.HashPassword(refresh); err != nil {
        return "", err
    }
	
    err = userCrud.UpdateRefreshTokenByUserId(ctx, tx, userId, refreshHashed, issuedAt)
	refresh = base64.StdEncoding.EncodeToString([]byte(refresh))
    
    return refresh, err
}

func UpdateToken(ctx context.Context, tx pgx.Tx, ip, jwt_token, refreshBase64 string) (string, string, error) {
    claims, err := ParseClaims(jwt_token)
	if err != nil {
		return "", "", err
	}
	refreshBytes, err := base64.StdEncoding.DecodeString(refreshBase64)
	if err != nil {
		return "", "", err
	}
	refresh := string(refreshBytes)

	user := userCrud.GetUserById(ctx, tx, claims.UserId)

	if !passwordUtils.ValidatePassword(refresh, user.Refresh) {
		return "", "", errors.New("Invalid refresh token")
	}

	if !user.RefreshIssuedAt.Equal(claims.IssuedAt.Time) {
		return "", "", errors.New("Refresh token was issued with other JWT")
	}

	if ip != claims.IP {
		go mailUtils.SendWarning(user.Email, ip)
	}

	return CreateToken(ctx, tx, user.Email, ip)
}

func JWTMiddleware(ctx *gin.Context) {
	jwt_header := ctx.GetHeader("Authorization")
	token, is_bearer := strings.CutPrefix(jwt_header, "Bearer ")

	if !is_bearer {
		ctx.String(http.StatusUnauthorized, "JWT is required")
        ctx.Abort()
		return
	}

    claims, err := ParseClaims(token)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
        ctx.Abort()
        return
	}

	tx := db.GetReadOnlyTransaction(ctx)
	defer tx.Rollback(ctx)
	user := userCrud.GetUserById(ctx, tx, claims.UserId)
	if user == nil || !user.RefreshIssuedAt.Equal(claims.IssuedAt.Time) {
		ctx.String(http.StatusBadRequest, "JWT is expired")
		ctx.Abort()
	}

	ctx.Set("claims", claims)
}

func ParseClaims(token string) (*Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(token, &claims, secretKeyFunc, jwt.WithExpirationRequired())
	return &claims, err
}

func secretKeyFunc(token *jwt.Token) (interface{}, error) {
    return SECRET_KEY, nil
}
