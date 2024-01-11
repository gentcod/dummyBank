package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createTransferRequest struct {
	SenderID string `json:"sender_id" binding:"required,uuid"`
	RecipientID string `json:"recipient_id" binding:"required,uuid"`
	Amount    int64  `json:"amount" binding:"required,min=1"`
	Currency  string    `json:"currency" binding:"required,currency"`
}

func(server *Server) getTransferById(ctx *gin.Context) {
	var req getEntityByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := server.store.GetTransfer(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

func (server *Server) getTransfers(ctx *gin.Context) {
	var req pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetTransfersParams{
		Limit: req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	entries, err := server.store.GetTransfers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

func (server *Server) createTransferTx(ctx *gin.Context) {
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validateAccount(ctx, uuid.MustParse(req.SenderID),req.Currency) {
		return
	}

	if !server.validateAccount(ctx, uuid.MustParse(req.RecipientID),req.Currency) {
		return
	}

	transfer, err := server.store.TransferTx(ctx, db.TransferTxParams{
		SenderID: uuid.MustParse(req.SenderID),
		RecipientID: uuid.MustParse(req.RecipientID),
		Amount: req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) validateAccount(ctx *gin.Context, accountId uuid.UUID, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %v currency mismatch: %v vs %v", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}