package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/SkiFoxi/Simple_bank/util"
	"github.com/stretchr/testify/require"
)

// createRandomAccount создает тестовый аккаунт и не будет работать как обычные тесты т.к. в названии нет префикса Test
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg) //testQueries объявлена в main_test
	require.NoError(t, err)                                              //t это t *testing.T Если ошибка равна 0 то все пройдет, если нет, то завалит тест
	require.NotEmpty(t, account)                                         //проверяет не является ли account пустым

	require.Equal(t, arg.Owner, account.Owner) // Проверка соответствуют ли аргументы заданным в тестовой таблице
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID) //Проверка генерируется ли ID для нового аккаунта
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner) //т.к. это передача данных из аккаунта 1 на аккаунт 2 то они должны быть равны
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
	//Это проверка на время создание, если оно в пределах заданного времени time.Second, то тест пройдет
}

func TestUpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	arg := UpdateAccountParams {
		ID: acc1.ID,       
		Balance: util.RandomMoney(), 
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner) //т.к. это передача данных из аккаунта 1 на аккаунт 2 то они должны быть равны
	require.Equal(t, arg.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)	
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)	
	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	
	require.NoError(t, err)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)  //Ожидаем ошибку - так как аккаунт должен быть удален
	require.EqualError(t, err, sql.ErrNoRows.Error()) //Проверяем конкретный тип ошибки - это должна быть ошибка "нет строк"
	require.Empty(t, acc2) //Дополнительная проверка, что возвращенный объект аккаунта пустой
}

func TestListAccounts(t *testing.T) {
	for i := 0; i<10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams {
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
} 