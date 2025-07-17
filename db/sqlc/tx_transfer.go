package db

import "context"

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// * transfertx oerfirms a money transfer from one account to another
// * it creates a transfer record, add account entries, update accounts balance within a single database transation
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	// Ensure consistent locking order
	var fromAccountID, toAccountID int64
	if arg.FromAccountId < arg.ToAccountId {
		fromAccountID = arg.FromAccountId
		toAccountID = arg.ToAccountId
	} else {
		fromAccountID = arg.ToAccountId
		toAccountID = arg.FromAccountId
	}

	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.GetAccountForUpdate(ctx, fromAccountID)
		if err != nil {
			return err
		}
		_, err = q.GetAccountForUpdate(ctx, toAccountID)
		if err != nil {
			return err
		}

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountId,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountId,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
