package gapi

import (
	"context"
	"database/sql"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
import "github.com/auronvila/simple-bank/pb/account"

func (server *Server) GetAccountById(ctx context.Context, req *account.GetAccountByIdRequest) (*account.GetAccountByIdResponse, error) {
	_, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user is not authenticated to perform the request")
	}

	getAccount, err := server.store.GetAccount(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.Unauthenticated, "could not find the account with the specified id")
		}
		return nil, status.Errorf(codes.Unauthenticated, "could not find the account")
	}

	convertedAccount := convertAccountForGet(getAccount)

	return convertedAccount, nil
}
