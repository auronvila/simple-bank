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

func convertTransfer(transfer db.Transfer) *acPb.Transfer {
	return &acPb.Transfer{
		Id:            transfer.ID,
		FromAccountId: transfer.FromAccountID,
		ToAccountId:   transfer.ToAccountID,
		Amount:        transfer.Amount,
		CreatedAt:     timestamppb.New(transfer.CreatedAt),
	}
}

func convertAccountForTransaction(account db.Account) *acPb.Account {
	return &acPb.Account{
		Id:        account.ID,
		Username:  account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: timestamppb.New(account.CreatedAt),
	}
}

func convertEntry(entry db.Entry) *acPb.Entry {
	return &acPb.Entry{
		Id:        entry.ID,
		AccountId: entry.AccountID,
		Amount:    entry.Amount,
		CreatedAt: timestamppb.New(entry.CreatedAt),
	}
}

func convertTransferTxRequest(transferObj db.TransferTxResult) *acPb.CreateTransferResponse {
	return &acPb.CreateTransferResponse{
		Transfer:    convertTransfer(transferObj.Transfer),
		FromAccount: convertAccountForTransaction(transferObj.FromAccount),
		ToAccount:   convertAccountForTransaction(transferObj.ToAccount),
		FromEntry:   convertEntry(transferObj.FromEntry),
		ToEntry:     convertEntry(transferObj.ToEntry),
	}
}
