package db

import (
	"context"
	"strconv"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createRandomTransaction(t *testing.T) Transaction {
	account := CreateRandomAccount(t)

	rm, _ := faker.RandomInt(1, 12, 1)
	rm_ := int64(rm[0])
	randomMonth := strconv.FormatInt(rm_, 10)

	rd, _ := faker.RandomInt(1, 27, 1)
	rd_ := int64(rd[0])
	randomDay := strconv.FormatInt(rd_, 10)

	rt, _ := faker.RandomInt(0, 100, 1)
	randomTransaction := float64(rt[0]) - 50.0

	arg := CreateTransactionParams{
		AccountID:   account.ID,
		Date:        randomMonth + "/" + randomDay,
		Transaction: randomTransaction,
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, arg.AccountID, transaction.AccountID)
	require.Equal(t, arg.Date, transaction.Date)
	require.Equal(t, arg.Transaction, transaction.Transaction)

	require.NotZero(t, transaction.ID)
	return transaction
}

func TestCreateTransaction(t *testing.T) {
	createRandomTransaction(t)
}

func TestGetTransaction(t *testing.T) {
	transaction1 := createRandomTransaction(t)
	transaction2, err := testQueries.GetTransaction(context.Background(), transaction1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transaction2)

	require.Equal(t, transaction1.ID, transaction2.ID)
	require.Equal(t, transaction1.AccountID, transaction2.AccountID)
	require.Equal(t, transaction1.Date, transaction2.Date)
	require.Equal(t, transaction1.Transaction, transaction2.Transaction)
}

func TestListTransactions(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransaction(t)
	}

	arg := ListTransactionsParams{
		Limit:  5,
		Offset: 5,
	}

	transactions, err := testQueries.ListTransactions(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transactions, 5)

	for _, tran := range transactions {
		require.NotEmpty(t, tran)
	}

}
