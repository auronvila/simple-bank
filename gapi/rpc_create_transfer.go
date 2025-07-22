package gapi

import (
	"context"
	"errors"
	"fmt"
	db "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/auronvila/simple-bank/pb/account"
	"github.com/auronvila/simple-bank/validator"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateTransfer(ctx context.Context, req *account.CreateTransferRequest) (*account.CreateTransferResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user is not authenticated to perform the request")
	}

	violations := ValidateCreateTransferReq(req)
	if violations != nil {
		return nil, InvalidArgument(violations)
	}

	fromAccount, err := server.validAccount(ctx, req.FromAccountId, req.Currency)
	if err != nil {
		return nil, err
	}

	if fromAccount.Owner != authPayload.Username {
		return nil, status.Errorf(codes.Canceled, "from account does not belong to the authenticated user")
	}

	_, err = server.validAccount(ctx, req.ToAccountId, req.Currency)
	if err != nil {
		return nil, err
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Amount:        req.Amount,
	}

	transferRes, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, "error in creating a transfer")
	}

	resObj := convertTransferTxRequest(transferRes)

	return resObj, nil
}

func (server *Server) validAccount(ctx context.Context, accountId int64, currency string) (*db.Account, error) {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "account with ID %d not found", accountId)
		}
		return nil, status.Errorf(codes.Internal, "failed to get account: %v", err)
	}

	if account.Currency != currency {
		err := fmt.Errorf("currency mismatch: account [%d] is %s, requested %s", account.ID, account.Currency, currency)
		log.Err(err).Msg("invalid currency")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &account, nil
}

func ValidateCreateTransferReq(req *account.CreateTransferRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, FieldViolation("currency", err))
	}

	return violations
}
