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

//Server serves HTTP requests for our banking service
type Server struct {
	config util.Config
	store db.Store
	tokenGenerator token.Generator
	router *gin.Engine
}

//NewServer creates a new HTTP server amd setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenGenerator, err := token.NewPasetoGenerator(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize token generator: %v", err)
	}

	server := &Server{
		config: config,
		store: store,
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

	v1Routes := router.Group("/api/v1")

	v1Routes.POST("/user/signup", server.createUser)	
	v1Routes.POST("/user/login", server.loginUser)
	v1Routes.PATCH("/user/update", server.updateUser)
	v1Routes.POST("/user/session/refresh", server.refreshSession)

	authRoutes := v1Routes.Group("/").Use((authMiddleware(server.tokenGenerator)))

	authRoutes.POST("/account", server.createAccount)
	authRoutes.PATCH("account", server.updateAccount)
	authRoutes.GET("/account", server.getAccounts)
	authRoutes.GET("/account/:id", server.getAccountById)

	authRoutes.POST("/transfer", server.createTransferTx)
	authRoutes.GET("/transfer", server.getTransfers)
	authRoutes.GET("/transfer/:id", server.getTransferById)

	server.router = router
}

// Start runs HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func apiErrorResponse(message string) gin.H {
	return gin.H{
		"status": "error",
		"message": message,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}