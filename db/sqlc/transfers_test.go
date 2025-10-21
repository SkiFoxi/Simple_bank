package db

import (
	"context"
	"testing"
	"time"

	"github.com/SkiFoxi/Simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfers(t *testing.T, account1 Account, account2 Account) Transfer {
	arg := CreateTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: util.RandomInt(1, 10),
	}

	transfer, err := testQueries.CreateTransfers(context.Background(), arg)
	require.NoError(t, err)                                              //t это t *testing.T Если ошибка равна 0 то все пройдет, если нет, то завалит тест
	require.NotEmpty(t, transfer)                                         //проверяет не является ли account пустым

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID) // Проверка соответствуют ли аргументы заданным в тестовой таблице
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID) //Проверка генерируется ли ID для нового аккаунта
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfers(t, account1, account2)
}

func TestGetTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfers(t, account1, account2)
	transfer2, err := testQueries.GetTransfers(context.Background(), transfer1.ID)

	require.NoError(t, err)                                             
	require.NotEmpty(t, transfer2)
	
	require.Equal(t, transfer1.ID, transfer2.ID) // Проверка соответствуют ли аргументы заданным в тестовой таблице
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount) // Проверка соответствуют ли аргументы заданным в тестовой таблице

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)	
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	for i := 0; i < 5; i++{
		createRandomTransfers(t, account1, account2)
		createRandomTransfers(t, account2, account1)
	}
	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Limit: 5,
		Offset: 0,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)	

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)	
		require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
		require.Equal(t, arg.ToAccountID, transfer.ToAccountID)	
	}
}

