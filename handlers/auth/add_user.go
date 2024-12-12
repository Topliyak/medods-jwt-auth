package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/medods-jwt-auth/db"
	"github.com/medods-jwt-auth/utils/password"
	userCrud "github.com/medods-jwt-auth/utils/crud/user"
)

func addUser(ctx *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var err error
	if err = ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var passwordHashed string
	if passwordHashed, err = password.HashPassword(body.Password); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	tx := db.GetReadWriteTransaction(ctx)
	defer tx.Rollback(ctx)

	err = userCrud.CreateUser(ctx, tx, body.Email, passwordHashed)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusCreated)
	tx.Commit(ctx)
}
