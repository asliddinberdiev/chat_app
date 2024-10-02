package router

import (
	"github.com/asliddinberdiev/chat_app/internal/user"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userH *user.Handler) {
	r = gin.Default()

	r.POST("/signup", userH.Create)
	r.POST("/login", userH.Login)
	r.GET("/logout", userH.Logout)
}

func Start(addr string) error {
	return r.Run(addr)
}
