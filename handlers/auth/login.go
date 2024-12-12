package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/medods-jwt-auth/db"
	"github.com/medods-jwt-auth/utils/auth"
	"github.com/medods-jwt-auth/utils/token"
)

func login(ctx *gin.Context) {
	type Credentials struct {
		Email string `json:"email" bind:"required"`
		Password string `json:"password" bind:"required"`
	}

	var creds Credentials

	if err := ctx.ShouldBindBodyWithJSON(&creds); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	tx := db.GetReadWriteTransaction(ctx)
	defer tx.Rollback(ctx)

	if !auth.CheckCredentials(ctx, tx, creds.Email, creds.Password) {
		ctx.String(http.StatusBadRequest, "Invalid credentials")
		return
	}

	jwt, refresh, err := token.CreateToken(ctx, tx, creds.Email, ctx.ClientIP())
	
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"jwt": jwt,
		"refresh": refresh,
	})

	tx.Commit(ctx)
}
