package api

import (
	simplebank "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createTransferReq struct {
	FromAccountId int64 `json:"from_account_id" binding:"required"`
	ToAccountId   int64 `json:"to_account_id"  binding:"required"`
	Amount        int64 `json:"amount"  binding:"required"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var reqData createTransferReq

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := simplebank.TransferTxParams{
		FromAccountId: reqData.FromAccountId,
		ToAccountId:   reqData.ToAccountId,
		Amount:        reqData.Amount,
	}

	transferRes, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transferRes)
}
