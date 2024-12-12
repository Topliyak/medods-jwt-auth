package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/medods-jwt-auth/db"
	"github.com/medods-jwt-auth/utils/token"
)

func refresh(ctx *gin.Context) {
	var reqBody struct {
		JWT string `json:"jwt" bind:"required"`
		Refresh string `json:"refresh" bind:"required"`
	}

	if err := ctx.ShouldBindBodyWithJSON(&reqBody); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	tx := db.GetReadWriteTransaction(ctx)
	defer tx.Rollback(ctx)

	jwt, refresh, err := token.UpdateToken(ctx, tx, ctx.ClientIP(), reqBody.JWT, reqBody.Refresh)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"jwt": jwt,
		"refresh": refresh,
	})

	tx.Commit(ctx)
}
