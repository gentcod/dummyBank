package api

import (
	"database/sql"
	"errors"
	"time"

	"net/http"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

type updateAccountRequest struct {
	AccountID string `json:"account_id" binding:"required,uuid"`
	Balance   int64  `json:"balance" binding:"required,min=0"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		ID:       uuid.New(),
		Owner:    authPayload.UserID,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, handlerResponse(ApiResponse[error]{
					StatusCode: http.StatusForbidden,
					Message:    err.Error(),
					Data:       nil,
				}))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	userAccount := getUserAccount(authPayload.Username, account)

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[UserAccount]{
		StatusCode: http.StatusOK,
		Message:    "account has been created successfully",
		Data:       userAccount,
	}))
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	arg := db.UpdateAccountParams{
		ID:        uuid.MustParse(req.AccountID),
		Balance:   req.Balance,
		UpdatedAt: time.Now(),
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[db.Account]{
		StatusCode: http.StatusOK,
		Message:    "account has been updated successfully",
		Data:       account,
	}))
}

func (server *Server) getAccountById(ctx *gin.Context) {
	var req getEntityByIdUUIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	account, err := server.store.GetAccount(ctx, uuid.MustParse(req.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, handlerResponse(ApiResponse[error]{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.ID {
		ctx.JSON(http.StatusUnauthorized, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusUnauthorized,
			Message:    errors.New("account doesn't belong to the authenticated user").Error(),
			Data:       nil,
		}))
		return
	}

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[db.Account]{
		StatusCode: http.StatusOK,
		Message:    "account has been fetched successfully",
		Data:       account,
	}))
}

func (server *Server) getAccounts(ctx *gin.Context) {
	var req pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetAccountsParams{
		Owner:  authPayload.UserID,
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	accounts, err := server.store.GetAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[[]db.Account]{
		StatusCode: http.StatusOK,
		Message:    "accounts have been fetched successfully",
		Data:       accounts,
	}))
}

func getUserAccount(username string, account db.Account) UserAccount {
	return UserAccount{
		Username:  username,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}
