package router

import (
	"github.com/asliddinberdiev/chat_app/internal/user"
	"github.com/asliddinberdiev/chat_app/internal/ws"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userH *user.Handler, wsH *ws.Handler) {
	r = gin.Default()

	r.POST("/signup", userH.Create)
	r.POST("/login", userH.Login)
	r.GET("/logout", userH.Logout)

	r.POST("/ws/create-room", wsH.CreateRoom)
	r.GET("/ws/join-room/:room_id", wsH.JoinRoom)
	r.GET("/ws/get-rooms", wsH.GetRooms)
	r.GET("/ws/get-clients/:room_id", wsH.GetClients)
}

func Start(addr string) error {
	return r.Run(addr)
}
