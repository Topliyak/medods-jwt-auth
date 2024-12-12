package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(r *gin.RouterGroup) {
	r.POST("/login", login)
	r.PUT("/refresh", refresh)
	r.POST("/user", addUser)
}
