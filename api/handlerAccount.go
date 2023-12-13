package api

import (
	"time"

	"net/http"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createAccounRequest struct {
	Owner     string    `json:"owner" binding:"required"`
	Currency  string    `json:"currency" binding:"required,oneof=USD EUR"`
}

func(server *Server) createAccount(ctx *gin.Context) {
	var req createAccounRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.CreateAccountParams{
		ID: uuid.New(),
		Owner: req.Owner,
		Balance: 0,
		Currency: req.Currency,
		UpdatedAt: time.Now().UTC(),
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func(server *Server) getAllAccounts(ctx *gin.Context) {
	arg := db.GetAllAccountsParams{
		Limit: 20,
		Offset: 20,
	}

	accounts, err := server.store.GetAllAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, accounts)
}