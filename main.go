package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/medods-jwt-auth/config"
	"github.com/medods-jwt-auth/db"
	"github.com/medods-jwt-auth/handlers"
	"github.com/medods-jwt-auth/utils/mail"
)

func main() {
	db.Init()
	mail.Init()

	r := gin.Default()
    handlers.RegisterHandlers(&r.RouterGroup)

	r.Use(gin.BasicAuth(gin.Accounts{
		"admin": "123",
		"client": "111",
	}))

	r.GET("", func(ctx *gin.Context) {
		user, _ := ctx.MustGet(gin.AuthUserKey).(string)
		ctx.String(200, user)
	})

    addr := fmt.Sprintf("%s:%d", config.SERVICE_HOST, config.SERVICE_PORT)
    r.Run(addr)
}
