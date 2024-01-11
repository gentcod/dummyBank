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

//Pagination is used for setting limit and offset for api request to the database
type pagination struct {
	PageId int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

//GetEntityByIdRequest is used to set binding request for uri using uuid 
type getEntityByIdUUIDRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

//GetEntityByIdRequest is used to set binding request for uri using uuid 
type getEntityByIdRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

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

	router.POST("/users", server.createUser)	
	router.POST("/users/login", server.loginUser)

	router.PATCH("/users", server.updateUser)
	// router.GET("/users", server.getUsers)67

	server.router = router
}

// Start runs HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}