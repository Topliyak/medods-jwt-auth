package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/medods-jwt-auth/handlers/auth"
	jwtsample "github.com/medods-jwt-auth/handlers/jwt-sample"
)

func RegisterHandlers(r *gin.RouterGroup) {
	auth.RegisterHandlers(r.Group("/auth"))
	jwtsample.RegisterHandlers(r.Group("/jwt-sample"))
}
