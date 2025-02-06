package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TODO: Implement authorization for getting transfers
// TODO: Create transfer response struct

type createTransferRequest struct {
	SenderID    string `json:"sender_id" binding:"required,uuid"`
	RecipientID string `json:"recipient_id" binding:"required,uuid"`
	Amount      int64  `json:"amount" binding:"required,min=1"`
	Currency    string `json:"currency" binding:"required,currency"`
}

type RecipientAccountResponse struct {
	ID       uuid.UUID `json:"id"`
	Currency string    `json:"currency"`
	Owner    uuid.UUID `json:"owner"`
}

type createTransferResponse struct {
	Transfer         db.Transfer              `json:"transfer"`
	SenderAccount    db.Account               `json:"sender_account"`
	RecipientAccount RecipientAccountResponse `json:"recipient_account"`
	SenderEntry      db.Entry                 `json:"sender_entry"`
}

func (server *Server) getTransferById(ctx *gin.Context) {
	var req getEntityByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	transfer, err := server.store.GetTransfer(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, handlerResponse(ApiResponse[error]{
				statusCode: http.StatusNotFound,
				message:    err.Error(),
				data:       nil,
			}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[db.Transfer]{
		statusCode: http.StatusOK,
		message:    "transfer record has been fetched successfully",
		data:       transfer,
	}))
}

func (server *Server) getTransfers(ctx *gin.Context) {
	var req pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	arg := db.GetTransfersParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	transfers, err := server.store.GetTransfers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[[]db.Transfer]{
		statusCode: http.StatusOK,
		message:    "transfer records have been fetched successfully",
		data:       transfers,
	}))
}

func (server *Server) createTransferTx(ctx *gin.Context) {
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	senderAcc, valid := server.validateAccount(ctx, uuid.MustParse(req.SenderID), req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if senderAcc.Owner != authPayload.UserID {
		err := errors.New("sender account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusUnauthorized,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	recipientAcc, valid := server.validateAccount(ctx, uuid.MustParse(req.RecipientID), req.Currency)
	if !valid {
		return
	}

	transfer, err := server.store.TransferTx(ctx, db.TransferTxParams{
		SenderID:    senderAcc.ID,
		RecipientID: recipientAcc.ID,
		Amount:      req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := createTransferResponse{
		Transfer:      transfer.Transfer,
		SenderAccount: transfer.SenderAccount,
		RecipientAccount: RecipientAccountResponse{
			ID:       transfer.RecipientAccount.ID,
			Currency: transfer.RecipientAccount.Currency,
			Owner:    transfer.RecipientAccount.Owner,
		},
		SenderEntry: transfer.SenderEntry,
	}

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[createTransferResponse]{
		statusCode: http.StatusOK,
		message:    "transfer has been processed successfully",
		data:       resp,
	}))
}

func (server *Server) validateAccount(ctx *gin.Context, accountId uuid.UUID, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, handlerResponse(ApiResponse[error]{
				statusCode: http.StatusNotFound,
				message:    err.Error(),
				data:       nil,
			}))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
			data:       nil,
		}))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %v currency mismatch: %v vs %v", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
			data:       nil,
		}))
		return account, false
	}

	return account, true
}
