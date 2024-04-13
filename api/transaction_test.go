package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/bxcodec/faker/v4"
	mockdb "github.com/g-ton/stori-candidate/db/mock"
	db "github.com/g-ton/stori-candidate/db/sqlc"
	mockmail "github.com/g-ton/stori-candidate/mail/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransactionAPI(t *testing.T) {
	account := randomAccount()
	transaction := randomTransaction(account.ID)

	testCases := []struct {
		name          string
		transaction   db.Transaction
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			transaction: transaction,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				// Getting account (required to create a transaction)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				// Creating transaction
				store.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Times(1).
					Return(transaction, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransaction(t, recorder.Body, transaction)
			},
		},
		{
			name:        "Internal Server Error",
			transaction: transaction,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				// Getting account (required to create a transaction)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				// Creating transaction
				store.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Transaction{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Bad request",
			transaction: db.Transaction{
				ID: 1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// Very important
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()
			// Routes (Defining actions)
			server.Router.POST("/transactions", server.CreateTransaction)
			// Convert data to json
			jsonInputAccount, _ := json.Marshal(tc.transaction)
			request, err := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonInputAccount))
			require.NoError(t, err)
			// ServeHTTP is a method from gin
			server.Router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetTransactionAPI(t *testing.T) {
	account := randomAccount()
	transaction := randomTransaction(account.ID)

	testCases := []struct {
		name          string
		transactionID int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			transactionID: transaction.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetTransaction(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(transaction, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransaction(t, recorder.Body, transaction)
			},
		},
		{
			name:          "Not Found",
			transactionID: transaction.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetTransaction(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(db.Transaction{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:          "Internal Server Error",
			transactionID: transaction.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetTransaction(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(db.Transaction{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:          "Bad request",
			transactionID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetTransaction(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// Very important
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()
			// Routes (Defining actions)
			server.Router.GET("/transactions/:id", server.GetTransaction)

			url := fmt.Sprintf("/transactions/%d", tc.transactionID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListTransactionsAPI(t *testing.T) {
	// Creating more of one single transaction to be able to test a list of transactions
	account := randomAccount()
	transactions := []db.Transaction{
		randomTransaction(account.ID),
		randomTransaction(account.ID),
	}

	type args struct {
		pageID   int32
		pageSize int32
	}

	testCases := []struct {
		name          string
		args          args
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			args: struct {
				pageID   int32
				pageSize int32
			}{
				pageID:   1,
				pageSize: 5,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					ListTransactions(gomock.Any(), gomock.Any()).
					Times(1).
					Return(transactions, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransactions(t, recorder.Body, transactions)
			},
		},
		{
			name: "Bad request",
			args: struct {
				pageID   int32
				pageSize int32
			}{
				pageID:   0, // We set 0 to get a bad request: binding:"required,min=1"`
				pageSize: 5,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					ListTransactions(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			args: struct {
				pageID   int32
				pageSize int32
			}{
				pageID:   1,
				pageSize: 5,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					ListTransactions(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Transaction{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// Very important
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			// Routes (Defining actions)
			server.Router.GET("/transactions", server.ListTransactions)

			url := fmt.Sprintf("/transactions?page_id=%d&page_size=%d", tc.args.pageID, tc.args.pageSize)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetSummaryInfoByFile(t *testing.T) {
	os.Setenv("EMAIL_TEMPLATE_LEVEL", "1")
	correctPathCsvFile := "../files/txns.csv"
	incorrectPathCsvFile := "../files/txns.txt"

	testCases := []struct {
		name          string
		input         listTransactionsByFileRequest
		buildStubs    func(mail *mockmail.MockMail)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			input: listTransactionsByFileRequest{
				FilePath: correctPathCsvFile,
				Mails:    "jdamianjm@gmail.com,jdamianjn@gmail.com",
			},
			buildStubs: func(mail *mockmail.MockMail) {
				// build stubs
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Bad request",
			input: listTransactionsByFileRequest{
				FilePath: correctPathCsvFile,
			},
			buildStubs: func(mail *mockmail.MockMail) {
				// build stubs
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			input: listTransactionsByFileRequest{
				FilePath: incorrectPathCsvFile,
				Mails:    "jdamianjm@gmail.com,jdamianjn@gmail.com",
			},
			buildStubs: func(mail *mockmail.MockMail) {
				// build stubs
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal Server Error Sending Email",
			input: listTransactionsByFileRequest{
				FilePath: correctPathCsvFile,
				Mails:    "jdamianjm@gmail.com,jdamianjn@gmail.com",
			},
			buildStubs: func(mail *mockmail.MockMail) {
				// build stubs
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(fmt.Errorf("html/template: cannot Parse after Execute"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			// start test server and send request
			server := newTestServer(t, nil, mail)
			recorder := httptest.NewRecorder()

			// Routes (Defining actions)
			server.Router.POST("/sendSummaryInfoByFile", server.GetSummaryInfoByFile)
			// Convert data to json
			jsonInput, _ := json.Marshal(tc.input)
			request, err := http.NewRequest(http.MethodPost, "/sendSummaryInfoByFile", bytes.NewBuffer(jsonInput))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetSummaryInfoByDB(t *testing.T) {
	os.Setenv("EMAIL_TEMPLATE_LEVEL", "1")
	account := randomAccount()
	transaction := randomTransaction(account.ID)
	transactions := []db.Transaction{transaction}

	testCases := []struct {
		name          string
		input         listTransactionsByAccountRequest
		buildStubs    func(store *mockdb.MockStore, mail *mockmail.MockMail)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			input: listTransactionsByAccountRequest{
				AccountID: account.ID,
				Mails:     "jdamianjm@gmail.com,jdamianjn@gmail.com",
			},
			buildStubs: func(store *mockdb.MockStore, mail *mockmail.MockMail) {
				// build stubs for DB
				store.EXPECT().
					ListTransactionsByAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(transactions, nil)
				// build stubs for Mail
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Bad request",
			input: listTransactionsByAccountRequest{
				AccountID: account.ID,
			},
			buildStubs: func(store *mockdb.MockStore, mail *mockmail.MockMail) {
				// build stubs for DB
				store.EXPECT().
					ListTransactionsByAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
				// build stubs for Mail
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			input: listTransactionsByAccountRequest{
				AccountID: account.ID,
				Mails:     "jdamianjm@gmail.com,jdamianjn@gmail.com",
			},
			buildStubs: func(store *mockdb.MockStore, mail *mockmail.MockMail) {
				// build stubs for DB
				store.EXPECT().
					ListTransactionsByAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return([]db.Transaction{}, sql.ErrConnDone)
				// build stubs for Mail
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal Server Error Sending Email",
			input: listTransactionsByAccountRequest{
				AccountID: account.ID,
				Mails:     "jdamianjm@gmail.com,jdamianjn@gmail.com",
			},
			buildStubs: func(store *mockdb.MockStore, mail *mockmail.MockMail) {
				// build stubs for DB
				store.EXPECT().
					ListTransactionsByAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(transactions, nil)
				// build stubs for Mail
				mail.EXPECT().
					SendMail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(fmt.Errorf("html/template: cannot Parse after Execute"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// Very important
			defer ctrl.Finish()

			mail := mockmail.NewMockMail(ctrl)
			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store, mail)

			// start test server and send request
			server := newTestServer(t, store, mail)
			recorder := httptest.NewRecorder()

			// Routes (Defining actions)
			server.Router.POST("/sendSummaryInfoByDB", server.GetSummaryInfoByDB)
			// Convert data to json
			jsonInput, _ := json.Marshal(tc.input)
			request, err := http.NewRequest(http.MethodPost, "/sendSummaryInfoByDB", bytes.NewBuffer(jsonInput))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomTransaction(accountID int64) db.Transaction {
	randomIntA, _ := faker.RandomInt(0, 1000, 1)

	rm, _ := faker.RandomInt(1, 12, 1)
	rm_ := int64(rm[0])
	randomMonth := strconv.FormatInt(rm_, 10)

	rd, _ := faker.RandomInt(1, 27, 1)
	rd_ := int64(rd[0])
	randomDay := strconv.FormatInt(rd_, 10)

	rt, _ := faker.RandomInt(0, 100, 1)
	randomTransaction := float64(rt[0]) - 50.0

	return db.Transaction{
		ID:          int64(randomIntA[0]),
		AccountID:   accountID,
		Date:        randomMonth + "/" + randomDay,
		Transaction: randomTransaction,
	}
}

func requireBodyMatchTransaction(t *testing.T, body *bytes.Buffer, transaction db.Transaction) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTransaction db.Transaction
	err = json.Unmarshal(data, &gotTransaction)
	require.NoError(t, err)
	require.Equal(t, transaction, gotTransaction)
}

func requireBodyMatchTransactions(t *testing.T, body *bytes.Buffer, transactions []db.Transaction) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTransactions []db.Transaction
	err = json.Unmarshal(data, &gotTransactions)
	require.NoError(t, err)
	require.Equal(t, transactions, gotTransactions)
}
