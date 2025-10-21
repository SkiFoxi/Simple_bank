package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Хранилище которое будет представлять все функции
// для индивидуального выполнения запросов к БД
// Чтобы транзакции работали нам нужно расширить Queries который создан sqlc
// т.к. Queries может выполнять только одну операцию над таблицей,
// а в транзакции нам придется использовать несколько операций одновременно
// Этот метод называется композицией
type Store struct {
	*Queries
	db *sql.DB //для создания новой транзакции для БД
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db), //New функция созданная sqlc которая возвращает тип Queries
	}
}

// Создает объект Queries с транзакцией и вызывает функцию обратного вызова
// Потом фиксирует или откатывает транзакцию в зависимости от выданной ошибки
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) //BeginTx создает транзакцию а TxOptions{} настраивает уровень изоляции
	//  но нам нужен стандартный уровень поэтому пишем n
	if err != nil {
		return err
	}
	//Вызываем функцию New() с созданной транзакцией
	//Функция работает с tx из-за того что тип *Tx соответствует интерфейсу DBTX и DBTX реализует как тип DB так и TX
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
//Структура которая содержит все необходимые входные параметры для перевода денег между счетами
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

//Структура которая содержит результат транзакций
type TransferTxResult struct{
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account`
	ToAccount Account `json:"to_account"`	
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"` 
}

// Создаст новую запись о переводе, добавит новые записи по счетам и обновит баланс счетов за одну транзакцию
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfers(ctx, CreateTransfersParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,

		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		}else{
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)	
		}
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		} 
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
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
)(account1 Account, account2 Account, err error) {           //тут <-1
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil {
		return           //такое написание return будет правильным, он вернет переменные account1, account2, err , т.к. они прописаны тут 1->
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})
	
	return
}