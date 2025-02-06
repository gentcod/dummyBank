package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type refreshSessionRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshSessionResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) refreshSession(ctx *gin.Context) {
	var req refreshSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	refreshPayload, err := server.tokenGenerator.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, handlerResponse(ApiResponse[error]{
				StatusCode: http.StatusNotFound,
				Message:    "refresh token not found",
				Data:       nil,
			}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusUnauthorized,
			Message:    "blocked session",
			Data:       nil,
		}))
		return
	}

	if session.Username != refreshPayload.Username {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusUnauthorized,
			Message:    "incorrect session user",
			Data:       nil,
		}))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusUnauthorized,
			Message:    "refresh token mismatch",
			Data:       nil,
		}))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusUnauthorized,
			Message:    "session has expired. Login again",
			Data:       nil,
		}))
		return
	}

	user, err := server.store.GetUser(ctx, refreshPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, handlerResponse(ApiResponse[error]{
				StatusCode: http.StatusNotFound,
				Message:    "User not found",
				Data:       nil,
			}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	accessToken, accessPayload, err := server.tokenGenerator.CreateToken(refreshPayload.Username, user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	resp := refreshSessionResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[refreshSessionResponse]{
		StatusCode: http.StatusOK,
		Message:    "token has been refreshed successfully",
		Data:       resp,
	}))
}
