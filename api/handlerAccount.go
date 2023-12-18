package api

import (
	"database/sql"
	"time"

	"net/http"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createAccountRequest struct {
	Owner     string    `json:"owner" binding:"required"`
	Currency  string    `json:"currency" binding:"required,oneof=USD EUR"`
}

type updateAccountRequest struct {
	AccountID string `json:"account_id" binding:"required,uuid"`
	Balance   int64     `json:"balance" binding:"required,min=0"`
}

func(server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
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

func(server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID: uuid.MustParse(req.AccountID),
		Balance: req.Balance,
		UpdatedAt: time.Now(),
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func(server *Server) getAccountById(ctx *gin.Context) {
	var req getEntityByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, uuid.MustParse(req.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func(server *Server) getAccounts(ctx *gin.Context) {
	var req pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	accounts, err := server.store.GetAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}