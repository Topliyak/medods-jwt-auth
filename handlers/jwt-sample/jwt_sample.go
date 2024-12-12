package jwtsample

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/medods-jwt-auth/utils/token"
)

func RegisterHandlers(r *gin.RouterGroup) {
	r.Use(token.JWTMiddleware)
	r.GET("/accounts/me", getUserInfo)
}

func getUserInfo(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	ctx.JSON(http.StatusOK, claims)
}
