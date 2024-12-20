package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"wallet-service/internal/db"
	"wallet-service/internal/db/mocks"
)

func Test_PostWalletOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var tests = []struct {
		name        string
		requestBody []byte
		statusCode  int
		repoMock    func() *mocks.MockRepository
	}{
		{
			name: "Deposit success",
			requestBody: []byte(`{
				"walletId": "123e4567-e89b-12d3-a456-426614174000",
				"operationType": "DEPOSIT",
				"amount": 100
			}`),
			statusCode: http.StatusOK,
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)

				repo.EXPECT().DepositMoney("123e4567-e89b-12d3-a456-426614174000", int64(100)).Return(nil)

				return repo
			},
		},
		{
			name: "DepositMoney error",
			requestBody: []byte(`{
					"walletId": "123e4567-e89b-12d3-a456-426614174000",
					"operationType": "DEPOSIT",
					"amount": 100
				}`),
			statusCode: http.StatusInternalServerError,
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)

				repo.EXPECT().DepositMoney("123e4567-e89b-12d3-a456-426614174000", int64(100)).Return(fmt.Errorf("random error"))

				return repo
			},
		},
		{
			name: "DepositMoney no wallet id in request",
			requestBody: []byte(`{
				"walletId": "",
				"operationType": "DEPOSIT",
				"amount": 100
			}`),
			statusCode: http.StatusBadRequest,
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)

				return repo
			},
		},
		{
			name: "DepositMoney amount = 0",
			requestBody: []byte(`{
				"walletId": "123e4567-e89b-12d3-a456-426614174000",
				"operationType": "DEPOSIT",
				"amount": 0
			}`),
			statusCode: http.StatusBadRequest,
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)

				return repo
			},
		},
		{
			name: "body operation type not present",
			requestBody: []byte(`{
				"walletId": "123e4567-e89b-12d3-a456-426614174000",
				"operationType": "",
				"amount": 100
			}`),
			statusCode: http.StatusBadRequest,
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)

				return repo
			},
		},
		{
			name: "Withdraw success",
			requestBody: []byte(`{
				"walletId": "123e4567-e89b-12d3-a456-426614174000",
				"operationType": "WITHDRAW",
				"amount": 100
			}`),
			statusCode: http.StatusOK,
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)

				repo.EXPECT().WithdrawMoney("123e4567-e89b-12d3-a456-426614174000", int64(100)).Return(nil)

				return repo
			},
		},
		{
			name: "WithdrawtMoney error",
			requestBody: []byte(`{
					"walletId": "123e4567-e89b-12d3-a456-426614174000",
					"operationType": "WITHDRAW",
					"amount": 100
				}`),
			statusCode: http.StatusInternalServerError,
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)

				repo.EXPECT().WithdrawMoney("123e4567-e89b-12d3-a456-426614174000", int64(100)).Return(fmt.Errorf("random error"))

				return repo
			},
		},
		{
			name: "WithdrawMoney no wallet id in request",
			requestBody: []byte(`{
				"walletId": "",
				"operationType": "WITHDRAW",
				"amount": 100
			}`),
			statusCode: http.StatusBadRequest,
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)

				return repo
			},
		},
		{
			name: "WithdrawMoney amount = 0",
			requestBody: []byte(`{
				"walletId": "123e4567-e89b-12d3-a456-426614174000",
				"operationType": "WITHDRAW",
				"amount": 0
			}`),
			statusCode: http.StatusBadRequest,
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)

				return repo
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repoMock := test.repoMock()
			handlerMocked := NewWalletHandler(repoMock)

			path := "/wallet"

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST(path, handlerMocked.PostWalletOperation)

			req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(test.requestBody))
			if err != nil {
				t.Errorf("http.NewRequest: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.statusCode, resp.Code)
		})
	}
}

func Test_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var tests = []struct {
		name         string
		walletUUID   string
		expectedBody []byte
		statusCode   int
		repoMock     func() *mocks.MockRepository
	}{
		{
			name:       "Get Balance success",
			walletUUID: "123e4567-e89b-12d3-a456-426614174000",
			statusCode: http.StatusOK,
			expectedBody: []byte(`{
				"balance": 500,
				"walletId": "123e4567-e89b-12d3-a456-426614174000"
			}`),
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)
				repo.EXPECT().GetBalance("123e4567-e89b-12d3-a456-426614174000").Return(int64(500), nil)
				return repo
			},
		},
		{
			name:       "Wallet not found",
			walletUUID: "123e4567-e89b-12d3-a456-426614174000",
			statusCode: http.StatusNotFound,
			expectedBody: []byte(`{
				"error": "Wallet not found"
			}`),
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)
				repo.EXPECT().GetBalance("123e4567-e89b-12d3-a456-426614174000").Return(int64(0), db.ErrWalletNotFound)
				return repo
			},
		},
		{
			name:       "Repository error",
			walletUUID: "123e4567-e89b-12d3-a456-426614174000",
			statusCode: http.StatusInternalServerError,
			expectedBody: []byte(`{
				"error": "Failed to fetch balance"
			}`),
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)
				repo.EXPECT().GetBalance("123e4567-e89b-12d3-a456-426614174000").Return(int64(0), fmt.Errorf("random error"))
				return repo
			},
		},
		{
			name:       "Missing walletUUID parameter",
			walletUUID: "",
			statusCode: http.StatusBadRequest,
			expectedBody: []byte(`{
				"error": "Missing walletUUID parameter"
			}`),
			repoMock: func() *mocks.MockRepository {
				repo := mocks.NewMockRepository(ctrl)
				return repo
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repoMock := test.repoMock()
			handlerMocked := NewWalletHandler(repoMock)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			if test.walletUUID == "" {
				router.GET("/wallets", handlerMocked.GetBalance)
			} else {
				router.GET("/wallets/:walletUUID", handlerMocked.GetBalance)
			}

			url := "/wallets"
			if test.walletUUID != "" {
				url = fmt.Sprintf("/wallets/%s", test.walletUUID)
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Errorf("http.NewRequest: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.statusCode, resp.Code)

			assert.JSONEq(t, string(test.expectedBody), resp.Body.String())
		})
	}
}
