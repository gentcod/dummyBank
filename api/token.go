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
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	refreshPayload, err := server.tokenGenerator.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusUnauthorized,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, handlerResponse(ApiResponse[error]{
				statusCode: http.StatusNotFound,
				message:    "refresh token not found",
				data:       nil,
			}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusUnauthorized,
			message:    "blocked session",
			data:       nil,
		}))
		return
	}

	if session.Username != refreshPayload.Username {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusUnauthorized,
			message:    "incorrect session user",
			data:       nil,
		}))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusUnauthorized,
			message:    "refresh token mismatch",
			data:       nil,
		}))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusUnauthorized,
			message:    "session has expired. Login again",
			data:       nil,
		}))
		return
	}

	user, err := server.store.GetUser(ctx, refreshPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, handlerResponse(ApiResponse[error]{
				statusCode: http.StatusNotFound,
				message:    "User not found",
				data:       nil,
			}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	accessToken, accessPayload, err := server.tokenGenerator.CreateToken(refreshPayload.Username, user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	resp := refreshSessionResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[refreshSessionResponse]{
		statusCode: http.StatusOK,
		message:    "token has been refreshed successfully",
		data:       resp,
	}))
}
