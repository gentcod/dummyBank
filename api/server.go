package api

import (
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//Server serves HTTP requests for our banking service
type Server struct {
	store db.Store
	router *gin.Engine
}

//Pagination is used for setting limit and offset for api request to the database
type pagination struct {
	PageId int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

//GetEntityByIdRequest is used to set binding request for uri using uuid 
type getEntityByIdRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

//NewServer creates a new HTTP server amd setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.PATCH("accounts", server.updateAccount)
	router.GET("/accounts", server.getAccounts)
	router.GET("/accounts/:id", server.getAccountById)

	router.POST("/entries", server.createEntry)
	router.GET("/entries", server.getEntries)
	router.GET("/entries/:id", server.getEntry)

	router.POST("/transfers", server.createTransferTx)
	router.GET("/transfers", server.getTransfers)
	router.GET("/transfers/:id", server.getTransferById)

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