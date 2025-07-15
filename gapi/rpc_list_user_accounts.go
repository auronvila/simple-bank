package gapi

import (
	"context"
	"database/sql"
	"errors"
	db "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/auronvila/simple-bank/pb/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ListUserAccounts(ctx context.Context, _ *account.ListUserAccountsRequest) (*account.ListUserAccountsResponse, error) {
	authorizedUser, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user is not authenticated to perform the request")
	}

	accountsParams := db.ListAccountsParams{
		Owner:  authorizedUser.Username,
		Limit:  99,
		Offset: 0,
	}

	accounts, err := server.store.ListAccounts(ctx, accountsParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "the user has no associated accounts")
		}
		return nil, status.Errorf(codes.Canceled, "error in listing user accounts")
	}

	var accountPointers []*account.Account
	for _, acc := range accounts {
		convertedAccount := convertAccount(acc)
		accountPointers = append(accountPointers, convertedAccount.Account)
	}

	return &account.ListUserAccountsResponse{Accounts: accountPointers}, nil
}
