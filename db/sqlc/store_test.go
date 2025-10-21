package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	t.Log(">> before:", account1.Balance, account2.Balance)
	n := 5 //звенья цепи транзакций 
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background() //создаем контекст в котором будет имя транзакции
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}
	existed := make(map[int] bool) // Маппа для k 
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)
		//Теперь проверим реально ли данные записались в БД
		_, err = store.GetTransfers(context.Background(), transfer.ID)
		require.NoError(t, err)

		//Теперь проверим записи в счете

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		//Теперь проверим получится ли получить запись об учетной записи из БД
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check accounts

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//check balance
		t.Log(">> after:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance //1000 - 800 значит списалось 200
		diff2 := toAccount.Balance - account2.Balance   // 1200 - 1000 значит прибавилось 200  diff1 и diff2 буду равны 200

		require.Equal(t, diff1, diff2)     //Проверка равны ли они
		require.True(t, diff1 > 0)         // проверяет действительно ли diff1 расчитался
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount... , n * amount  т.е. проверяет количество переводов если diff1 т.е. изменение баланса 200
		// а одна операция была на 50, значит 200 / 50, произошло 4 операции одновременно
		k := int(diff1 / amount)          // Расчитывает количество произведенных операций
		require.True(t, k >= 1 && k <= n) // Количество этих операций должно быть минимум 1 и не больше n
		require.NotContains(t, existed, k) //Проверка на то что в мапе existed нет k
		existed[k] = true  // добавляем в маппу 
	}

	//Проверяем обновленный счет 

//Проверяем обновленный счет 
updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
require.NoError(t, err)

updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)  // ← ИСПРАВЛЕНО!
require.NoError(t, err)	

fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

require.Equal(t, account1.Balance - int64(n)*amount, updatedAccount1.Balance)
require.Equal(t, account2.Balance + int64(n)*amount, updatedAccount2.Balance)
}


//Проверка на взаимоблокировку двух счетов в одной транзакции если менять счета через раз
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	t.Log(">> before:", account1.Balance, account2.Balance)
	n := 10 //звенья цепи транзакций 
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		 fromAccountID := account1.ID
		 toAccountID := account2.ID

		 if i % 2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		 }
		go func() {
			ctx := context.Background() //создаем контекст в котором будет имя транзакции
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	//Проверяем обновленный счет 

//Проверяем обновленный счет 
updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
require.NoError(t, err)

updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)  // ← ИСПРАВЛЕНО!
require.NoError(t, err)	

fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

require.Equal(t, account1.Balance, updatedAccount1.Balance)
require.Equal(t, account2.Balance, updatedAccount2.Balance)
}