package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries         //可以调用Quries的数据库操作(嵌套)
	db       *sql.DB //Store层的数据库操作
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db), //New(db)返回*Queries
	}
}

// execTx 执行数据库事务
// fn：接收*Queries参数并返回一个错误
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	//BeginTx starts a transaction.
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	//fn(q)执行了一系列的数据库操作，如查询、插入、更新等。err表示在执行数据库操作时是否发生了任何错误
	err = fn(q)
	if err != nil {
		//回滚数据(Rollback)是指在数据库事务中撤销对数据的更改操作，将数据恢复到事务开始之前的状态
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v,rb err:%v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTx的参数
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTx的结果
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account_id"`
	ToAccount   Account  `json:"to_account_id"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{} //

// 执行金钱交易
// 1.生成一个transfer记录：result.Transfer
// 2.添加2个账户
// 3.更新acounts的balance(余额)
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	// execTx 执行数据库事务
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		//debug
		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")
		// 1.生成一个transfer记录：result.Transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		// 2.添加2个账户
		fmt.Println(txName, "create entry1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//3.更新balance
		// fmt.Println(txName, "update account1")
		// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.FromAccountID,
		// 	Amount: -arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "update account2")
		// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.ToAccountID,
		// 	Amount: arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return
	}
	return
}
