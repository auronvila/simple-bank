package gapi

import (
	db "github.com/auronvila/simple-bank/db/sqlc"
	acPb "github.com/auronvila/simple-bank/pb/account"
	pb "github.com/auronvila/simple-bank/pb/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

func convertAccount(account db.Account) *acPb.CreateAccountResponse {
	return &acPb.CreateAccountResponse{
		Account: &acPb.Account{
			Id:        account.ID,
			Username:  account.Owner,
			Currency:  account.Currency,
			Balance:   account.Balance,
			CreatedAt: timestamppb.New(account.CreatedAt),
		},
	}
}

func convertAccountForUpdate(account db.Account) *acPb.UpdateAccountBalanceResponse {
	return &acPb.UpdateAccountBalanceResponse{
		Account: &acPb.Account{
			Id:        account.ID,
			Username:  account.Owner,
			Currency:  account.Currency,
			CreatedAt: timestamppb.New(account.CreatedAt),
			Balance:   account.Balance,
		},
	}
}

func convertAccountForGet(account db.Account) *acPb.GetAccountByIdResponse {
	return &acPb.GetAccountByIdResponse{Account: &acPb.Account{
		Id:        account.ID,
		Username:  account.Owner,
		Currency:  account.Currency,
		CreatedAt: timestamppb.New(account.CreatedAt),
		Balance:   account.Balance,
	}}
}
