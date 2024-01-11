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
	Username        string    `json:"username" binding:"required,alphanum"`
	FullName        string    `json:"full_name" binding:"required"`
	Email           string    `json:"email" binding:"required,email"`
	Password string    `json:"password" binding:"required,min=8"`
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
		Username: req.Username,
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

	userProfile := getUserProfile(user)

	ctx.JSON(http.StatusOK, userProfile)
}


type updateUserRequest struct {
	Username   string `json:"username" binding:"required,alphanum"`
	Password   string    `json:"password" binding:"required"`
	NewPassword   string    `json:"new_password" binding:"required"`
}

func(server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, valid := server.validateUser(ctx, req.Username, req.Password)
	if !valid {
		return
	}

	hashedNewPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID: user.ID,
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

//TODO: Implement multiple login methods; username or password
type loginUserRequest struct {
	Username        string    `json:"username" binding:"required,alphanum"`
	Password string    `json:"password" binding:"required,min=8"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User UserProfile `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, valid := server.validateUser(ctx, req.Username, req.Password)
	if !valid {
		return
	}

	accessToken, err := server.tokenGenerator.CreateToken(req.Username, server.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	userProfile := getUserProfile(user)

	resp := loginUserResponse{
		AccessToken: accessToken,
		User: userProfile,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) validateUser(ctx *gin.Context, username string, password string) (db.User, bool) {
	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return user, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return user, false
	}
	
	err = util.CheckPassword(password, user.HarshedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return user, false
	}

	return user, true
}

func getUserProfile(user db.User) UserProfile {
	return UserProfile{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}