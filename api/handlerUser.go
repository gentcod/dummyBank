package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Password string    `json:"password" binding:"required"`
	FullName        string    `json:"full_name" binding:"required"`
	Email           string    `json:"email" binding:"required"`
}

type updateUserRequest struct {
	UserID   string `json:"user_id" binding:"required,uuid"`
	Password   string    `json:"password" binding:"required"`
	NewPassword   string    `json:"new_password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		ID: uuid.New(),
		HarshedPassword: req.Password,
		FullName: req.FullName,
		Email: req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func(server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validateUser(ctx, uuid.MustParse(req.UserID), req.Password) {
		return
	}

}

func (server *Server) validateUser(ctx *gin.Context, userId uuid.UUID, harshedPassword string) bool {
	user, err := server.store.GetUserById(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if user.HarshedPassword != harshedPassword {
		err := fmt.Errorf("Wrong Password")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}