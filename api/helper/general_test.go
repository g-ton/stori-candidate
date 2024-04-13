package helper

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/bxcodec/faker/v4"
	db "github.com/g-ton/stori-candidate/db/sqlc"
	mockmail "github.com/g-ton/stori-candidate/mail/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

/*
As good habit implementing unit tests our code is the best way to ensure that
our features are working as expected without bugs nor regressions
*/

func TestGetSummaryInfo(t *testing.T) {
	transactions := []db.Transaction{}
	debitTotal, creditTotal := 0.0, 0.0
	debitNum, creditNum := 0, 0

	for i := 1; i <= 12; i++ {
		randomMonth := strconv.FormatInt(int64(i), 10)

		rd, _ := faker.RandomInt(1, 27, 1)
		rd_ := int64(rd[0])
		randomDay := strconv.FormatInt(rd_, 10)

		rt, _ := faker.RandomInt(0, 100, 1)
		randomTransaction := float64(rt[0]) - 50.0

		transactions = append(transactions, db.Transaction{
			ID:          int64(i),
			Date:        randomMonth + "/" + randomDay,
			Transaction: randomTransaction,
		})

		if randomTransaction < 0 { //Debit (-)
			debitTotal += randomTransaction
			debitNum++
		} else { //Credit (+)
			creditTotal += randomTransaction
			creditNum++
		}
	}

	gotTransactions := GetSummaryInfo(transactions)
	require.NotEmpty(t, gotTransactions)
	require.NotEmpty(t, gotTransactions.Months)

	require.Equal(t, debitTotal/float64(debitNum), gotTransactions.Data["avgDebit"])
	require.Equal(t, creditTotal/float64(creditNum), gotTransactions.Data["avgCredit"])
	require.Equal(t, creditTotal+debitTotal, gotTransactions.Data["totalBalance"])
}

func TestGetSummaryInfoMonths(t *testing.T) {
	transactions := []db.Transaction{
		{
			ID:          1,
			Date:        "1/10",
			Transaction: 20,
		},
		{
			ID:          2,
			Date:        "1/12",
			Transaction: -20,
		},
		{
			ID:          3,
			Date:        "12/10",
			Transaction: 20,
		},
		{
			ID:          4,
			Date:        "12/12",
			Transaction: -20,
		},
	}

	gotTransactions := GetSummaryInfo(transactions)
	require.NotEmpty(t, gotTransactions)
	require.NotEmpty(t, gotTransactions.Months)
	require.NotEmpty(t, gotTransactions.Data)

	require.Equal(t, 2, gotTransactions.Months["January"])
	require.Equal(t, 2, gotTransactions.Months["December"])

	// Checking with an empty transaction
	gotTransactions = GetSummaryInfo([]db.Transaction{})
	require.Empty(t, gotTransactions)
}

/*
cases for TestProcessFile

get transactions as expected
non existent file
*/
func TestProcessFile(t *testing.T) {
	testCases := []struct {
		name          string
		buildStubs    func() ([]db.Transaction, error)
		checkResponse func(t *testing.T, trans []db.Transaction, err error)
	}{
		{
			name: "OK",
			buildStubs: func() ([]db.Transaction, error) {
				return ProcessFile("../../files/txns.csv")
			},
			checkResponse: func(t *testing.T, trans []db.Transaction, err error) {
				// check response
				require.Nil(t, err)
				require.Equal(t, 4, len(trans))
			},
		},
		{
			name: "NonExistentFile",
			buildStubs: func() ([]db.Transaction, error) {
				return ProcessFile("../../files/txns_not_found.csv")
			},
			checkResponse: func(t *testing.T, trans []db.Transaction, err error) {
				// check response
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "no such file or directory")
				require.Nil(t, trans)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			// build stubs
			trans, err := tc.buildStubs()

			// check response
			tc.checkResponse(t, trans, err)
		})
	}
}

func TestProcessTemplateEmailForTransaction(t *testing.T) {
	tr := TransactionResult{
		Months: map[string]int{
			"January":  2,
			"February": 2,
		},
		Data: map[string]float64{
			"totalBalance": 10.0,
			"avgCredit":    10.0,
			"avgDebit":     -10.0,
		},
	}

	testCases := []struct {
		name          string
		input         TransactionResult
		buildStubs    func(mail *mockmail.MockMail)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name:  "OK",
			input: tr,
			buildStubs: func(mail *mockmail.MockMail) {
				// build stubs
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				// check response
				require.Nil(t, err)
			},
		},
		{
			name:  "Error Sending Mail",
			input: tr,
			buildStubs: func(mail *mockmail.MockMail) {
				// build stubs
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(fmt.Errorf("html/template: cannot Parse after Execute"))
			},
			checkResponse: func(t *testing.T, err error) {
				// check response
				require.NotNil(t, err)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// Very important
			defer ctrl.Finish()

			mail := mockmail.NewMockMail(ctrl)
			// build stubs
			tc.buildStubs(mail)

			err := ProcessTemplateEmailForTransaction(tc.input, []string{"jdamianjm@gmail.com"}, mail)

			// check response
			tc.checkResponse(t, err)
		})
	}
}
