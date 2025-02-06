package api

import (
	"database/sql"
	"net/http"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gin-gonic/gin"
)

func (server *Server) getEntry(ctx *gin.Context) {
	var req getEntityByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	entry, err := server.store.GetEntry(ctx, req.Id)
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

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[db.Entry]{
		StatusCode: http.StatusOK,
		Message:    "entry has been fetched successfully",
		Data:       entry,
	}))
}

func (server *Server) getEntries(ctx *gin.Context) {
	var req pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	arg := db.GetEntriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	entries, err := server.store.GetEntries(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handlerResponse(ApiResponse[error]{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		}))
		return
	}

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[[]db.Entry]{
		StatusCode: http.StatusOK,
		Message:    "entries have been fetched successfully",
		Data:       entries,
	}))
}
