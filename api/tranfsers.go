package api

import (
	"errors"
	"fmt"
	db "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/auronvila/simple-bank/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createTransferReq struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id"  binding:"required,min=1"`
	Amount        int64  `json:"amount"  binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var reqData createTransferReq

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, reqData.FromAccountId, reqData.Currency)
	if !valid {
		return
	}
	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, reqData.ToAccountId, reqData.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
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

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency missmatch %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
