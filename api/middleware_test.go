package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gentcod/DummyBank/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func addAuthorization(t *testing.T, request *http.Request, tokenGenerator token.Generator, authorizationType string, username string, duration time.Duration) {
	token, err := tokenGenerator.CreateToken(username, uuid.New(), duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%v %v", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenGenerator token.Generator)
		checkResponse func(t *testing.T, recoreder httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recoreder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoreder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
			},
			checkResponse: func(t *testing.T, recoreder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoreder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, "unsupported", "user", time.Minute)

			},
			checkResponse: func(t *testing.T, recoreder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoreder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, "", "user", time.Minute)

			},
			checkResponse: func(t *testing.T, recoreder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoreder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Generator) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recoreder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoreder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			testServer := testServerInit(t)

			authPath := "/auth"
			testServer.server.router.GET(authPath, authMiddleware(testServer.server.tokenGenerator), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, testServer.server.tokenGenerator)
			testServer.server.router.ServeHTTP(testServer.recorder, request)
			tc.checkResponse(t, *testServer.recorder)
		})
	}
}
