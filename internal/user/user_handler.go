package user

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) Create(ctx *gin.Context) {
	var input CreateUserReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.Create(ctx.Request.Context(), &input)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			log.Println("already used username or email: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "already used username or email"})
			return
		}
		log.Println("create user error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h *Handler) Login(ctx *gin.Context) {
	var input LoginReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.Login(ctx.Request.Context(), &input)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("not found user: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "email or password wrong"})
			return
		}

		log.Println("login user error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.SetCookie("jwt", user.AccessToken, 3600, "/", "0.0.0.0", false, true)

	ctx.JSON(http.StatusOK, gin.H{"id": user.ID, "username": user.Username})
}

func (h *Handler) Logout(ctx *gin.Context) {
	_, err := ctx.Cookie("jwt")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "not found"})
		return
	}
	ctx.SetCookie("jwt", "", -1, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
