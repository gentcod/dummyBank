package api

import (
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gin-gonic/gin"
)

//Server serves HTTP requests for our banking service
type Server struct {
	store *db.Store
	router *gin.Engine
}

//NewServer creates a new HTTP server amd setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.getAllAccounts)

	server.router = router
	return server
}

// Start runs HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}