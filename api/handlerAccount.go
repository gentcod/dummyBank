package api

import (
	"database/sql"
	"time"

	"net/http"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

//TODO: Implement User Validation to create and update an account

type createAccountRequest struct {
	UserID     string    `json:"user_id" binding:"required,uuid"`
	Currency  string    `json:"currency" binding:"required,currency`
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

	// user, err := server.store.GetUserById(ctx, uuid.MustParse(req.UserID))
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	arg := db.CreateAccountParams{
		ID: uuid.New(),
		Owner: uuid.MustParse(req.UserID),
		Balance: 0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name(){
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userAccount := UserAccount{
		ID: account.ID,
		// FullName: user.FullName,
		// Email: user.Email,
		Balance: account.Balance,
		Currency: account.Currency,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, userAccount)
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
			return
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