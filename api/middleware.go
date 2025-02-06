package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gentcod/DummyBank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func authMiddleware(tokenGenerator token.Generator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
				Data:       nil,
			}))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
				Data:       nil,
			}))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsuppoerted authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
				Data:       nil,
			}))
			return
		}

		accessToken := fields[1]
		payload, err := tokenGenerator.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
				Data:       nil,
			}))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
