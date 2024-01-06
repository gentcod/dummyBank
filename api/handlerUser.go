package api

import (
	"database/sql"
	"time"

	// "fmt"
	"net/http"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

//TODO: Implement password confirmation when updating account and password auth for getting account

type createUserRequest struct {
	FullName        string    `json:"full_name" binding:"required"`
	Email           string    `json:"email" binding:"required,email"`
	Password string    `json:"password" binding:"required,min=8"`
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

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		ID: uuid.New(),
		FullName: req.FullName,
		Email: req.Email,
		HarshedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name(){
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userProfile := UserProfile{
		ID: user.ID,
		FullName: user.FullName,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}

	ctx.JSON(http.StatusOK, userProfile)
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

	hashedNewPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID: uuid.MustParse(req.UserID),
		HarshedPassword: hashedNewPassword,
		PasswordChangedAt: time.Now().UTC(),
	}

	updatedUser, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func(server *Server) getUserById(ctx *gin.Context) {
	var req getEntityByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, uuid.MustParse(req.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func(server *Server) getUsers(ctx *gin.Context) {
	var req pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetUsersParams{
		Limit: req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	users, err := server.store.GetUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) validateUser(ctx *gin.Context, userId uuid.UUID, password string) bool {
	user, err := server.store.GetUserById(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	
	err = util.CheckPassword(password, user.HarshedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}