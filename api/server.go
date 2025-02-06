package api

import (
	"fmt"
	// "log"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/token"
	"github.com/gentcod/DummyBank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service
type Server struct {
	config         util.Config
	store          db.Store
	tokenGenerator token.Generator
	router         *gin.Engine
}

// NewServer creates a new HTTP server amd setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenGenerator, err := token.NewPasetoGenerator(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize token generator: %v", err)
	}

	server := &Server{
		config:         config,
		store:          store,
		tokenGenerator: tokenGenerator,
	}

	// Attach custom request validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(interceptor(server))

	v1Routes := router.Group("/api/v1")

	v1Routes.POST("/session/refresh", server.refreshSession)

	v1Routes.POST("/users/signup", server.createUser)
	v1Routes.POST("/users/login", server.loginUser)
	v1Routes.PATCH("/users/update", server.updateUser)

	authRoutes := v1Routes.Group("/").Use((authMiddleware(server.tokenGenerator)))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.PATCH("accounts", server.updateAccount)
	authRoutes.GET("/accounts", server.getAccounts)
	authRoutes.GET("/accounts/:id", server.getAccountById)

	authRoutes.POST("/transfers", server.createTransferTx)
	authRoutes.GET("/transfers", server.getTransfers)
	authRoutes.GET("/transfer/:id", server.getTransferById)

	authRoutes.GET("/transactions", server.getEntries)
	authRoutes.GET("/transactions/:id", server.getEntry)

	server.router = router
}

// Start runs HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

type ApiResponse[T any] struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

func handlerResponse[T any](resp ApiResponse[T]) gin.H {
	if resp.StatusCode < 300 {
		return apiSuccessResponse(resp)
	}

	return apiErrorResponse(resp)
}

func apiSuccessResponse[T any](resp ApiResponse[T]) gin.H {
	return gin.H{
		"status":     "success",
		"statusCode": resp.StatusCode,
		"message":    resp.Message,
		"data":       resp.Data,
	}
}

func apiErrorResponse[T any](resp ApiResponse[T]) gin.H {
	return gin.H{
		"status":     "error",
		"statusCode": resp.StatusCode,
		"message":    resp.Message,
		"data":       resp.Data,
	}
}

func handleInternalResponse(resp ApiResponse[any]) gin.H {
	return gin.H{
		"status":     "error",
		"statusCode": resp.StatusCode,
		"message":    "An unexpected error occured",
	}
}
