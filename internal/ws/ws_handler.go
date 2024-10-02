package ws

import (
	"log"
	"net/http"

	"github.com/asliddinberdiev/chat_app/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{hub: h}
}

type CreateRoomReq struct {
	Name string `json:"name"`
}

type CreateRoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(ctx *gin.Context) {
	var input CreateRoomReq

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := utils.UUID()
	if err != nil {
		log.Println("rrror generating room ID: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	for _, room := range h.hub.Rooms {
		if room.Name == input.Name {
			log.Println("room already exists with name: ", input.Name)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "room with this name already exists"})
			return
		}
	}

	resp := &CreateRoomRes{
		ID:   id,
		Name: input.Name,
	}

	h.hub.Rooms[resp.ID] = &Room{
		ID:      resp.ID,
		Name:    resp.Name,
		Clients: make(map[string]*Client),
	}

	ctx.JSON(http.StatusOK, resp)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := ctx.Param("room_id")
	clientID := ctx.Query("user_id")
	username := ctx.Query("username")

	cl := &Client{
		Conn:     conn,
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
		Message:  make(chan *Message, 10),
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(ctx *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, room := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   room.ID,
			Name: room.Name,
		})
	}

	ctx.JSON(http.StatusOK, rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(ctx *gin.Context) {
	var clients []ClientRes
	roomID := ctx.Param("room_id")

	if _, ok := h.hub.Rooms[roomID]; !ok {
		clients = make([]ClientRes, 0)
		ctx.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomID].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	ctx.JSON(http.StatusOK, clients)
}
