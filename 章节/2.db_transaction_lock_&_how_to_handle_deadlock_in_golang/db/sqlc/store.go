package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store 提供执行数据库查询和事务的所有功能
type Store struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx 在事务中执行一个函数
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// txOptions := &sql.TxOptions{
	// 	Isolation: sql.LevelDefault,
	// }
	// tx, err := store.db.BeginTx(ctx, txOptions) // BeginTx() 第二个参数可以设置隔离级别
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)
	err = fn(queries) // 执行事务发生的错误
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil { // 回滚发生的错误
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err // 回滚成功，返回事务错误
	}

	return tx.Commit() // 执行事务成功，提交事务
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx 将资金从一个帐户转到另一个帐户
// 它在数据库事务中创建转账，添加帐户条目并更新帐户余额
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// 当我们想从回调函数中获得结果时，通常使用闭包，因为回调函数本身不知道应返回的确切类型
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		// 转账对象
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		// 支出流水记录
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // 支出
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		// 收入流水记录
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // 收入
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get account 1")
		// move money out of account1
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get account 2")
		// move money into account2
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}
