package gapi

import (
	"context"
	"errors"
	db "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/auronvila/simple-bank/pb/account"
	"github.com/auronvila/simple-bank/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateAccountBalance(ctx context.Context, req *account.UpdateAccountBalanceRequest) (*account.UpdateAccountBalanceResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := ValidateUpdateAccountBalanceReq(req)
	if violations != nil {
		return nil, InvalidArgument(violations)
	}

	params := db.UpdateAccountBasedOnUsernameParams{
		Owner:    authPayload.Username,
		Balance:  req.Balance,
		Currency: req.Currency,
	}

	foundedAccount, err := server.store.UpdateAccountBasedOnUsername(ctx, params)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "no accounts were found")
		}

		return nil, status.Errorf(codes.Internal, "could not update the account %s", err)
	}

	convertedAccountForResponse := convertAccountForUpdate(foundedAccount)

	return convertedAccountForResponse, nil
}

func ValidateUpdateAccountBalanceReq(req *account.UpdateAccountBalanceRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, FieldViolation("currency", err))
	}

	if err := validator.ValidateBalance(req.GetBalance()); err != nil {
		violations = append(violations, FieldViolation("balance", err))
	}

	return violations
}
