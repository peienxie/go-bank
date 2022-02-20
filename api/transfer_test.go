package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/peienxie/go-bank/db/mock"
	db "github.com/peienxie/go-bank/db/sqlc"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransferAPI(t *testing.T) {
	amount := randomMoney()
	account1 := randomAccount()
	account2 := randomAccount()
	otherCurrencyAccount := randomAccount()

	currency := "USD"
	otherCurrency := "TWD"
	invalidCurrency := "zzz"
	account1.Currency = currency
	account2.Currency = currency
	otherCurrencyAccount.Currency = otherCurrency

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			"OK",
			gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        currency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(1).Return(account2, nil)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(1)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			"BadRequest request from account id invalid",
			gin.H{
				"from_account_id": -1,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        currency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(0)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			"BadRequest request to account id invalid",
			gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   -1,
				"amount":          amount,
				"currency":        currency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(0)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			"BadRequest request amount invalid",
			gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          0,
				"currency":        currency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(0)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			"BadRequest request currency invalid",
			gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        invalidCurrency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(0)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			"NotFound from acount not found",
			gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        currency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(0)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			"NotFound to acount not found",
			gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        currency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(1).Return(db.Account{}, sql.ErrNoRows)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			"InternalServerError get acount error",
			gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        currency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(1).Return(db.Account{}, sql.ErrConnDone)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			"BadRequest from account currency not match",
			gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        currency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(otherCurrencyAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(0)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			"BadRequest to account currency not match",
			gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        currency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(1).Return(otherCurrencyAccount, nil)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(0)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			"InternalServerError transfer error",
			gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        currency,
			},
			func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(1).Return(account2, nil)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			assert.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
