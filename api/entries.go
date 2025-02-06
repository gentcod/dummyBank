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
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	entry, err := server.store.GetEntry(ctx, req.Id)
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

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[db.Entry]{
		statusCode: http.StatusOK,
		message:    "entry has been fetched successfully",
		data:       entry,
	}))
}

func (server *Server) getEntries(ctx *gin.Context) {
	var req pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handlerResponse(ApiResponse[error]{
			statusCode: http.StatusBadRequest,
			message:    err.Error(),
			data:       nil,
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
			statusCode: http.StatusInternalServerError,
			message:    err.Error(),
			data:       nil,
		}))
		return
	}

	ctx.JSON(http.StatusOK, handlerResponse(ApiResponse[[]db.Entry]{
		statusCode: http.StatusOK,
		message:    "entries have been fetched successfully",
		data:       entries,
	}))
}
