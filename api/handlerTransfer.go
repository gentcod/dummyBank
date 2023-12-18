package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createTransferRequest struct {
	SenderID string `json:"sender_id" binding:"required,uuid"`
	RecipientID string `json:"recipient_id" binding:"required,uuid"`
	Amount    int64  `json:"amount" binding:"required,min=1"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	senderAccount, err := server.store.GetAccount(ctx, uuid.MustParse(req.SenderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if senderAccount.Balance < req.Amount {
		err = errors.New("insufficient balance")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	recipientAccount, err := server.store.GetAccount(ctx, uuid.MustParse(req.RecipientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateTransferParams{
		ID: uuid.New(),
		SenderID: senderAccount.ID,
		RecipientID: recipientAccount.ID,
		Amount: req.Amount,
	}

	transfer, err := server.store.CreateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//Create Entries
	_, err = server.store.CreateEntry(ctx, db.CreateEntryParams{
		ID: uuid.New(),
		AccountID: senderAccount.ID,
		Amount: -req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.CreateEntry(ctx, db.CreateEntryParams{
		ID: uuid.New(),
		AccountID: recipientAccount.ID,
		Amount: req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//Update Accounts
	_, err = server.store.UpdateAccount(ctx, db.UpdateAccountParams{
		ID: senderAccount.ID,
		Balance: senderAccount.Balance - req.Amount,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.UpdateAccount(ctx, db.UpdateAccountParams{
		ID: recipientAccount.ID,
		Balance: recipientAccount.Balance + req.Amount,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func(server *Server) getTransferById(ctx *gin.Context) {
	var req getEntityByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := server.store.GetTransfer(ctx, uuid.MustParse(req.Id))
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