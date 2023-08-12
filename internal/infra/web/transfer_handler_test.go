package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/entity/mock"
	"lucassantoss1701/bank/internal/infra/web"
	"lucassantoss1701/bank/internal/usecase"
	usecaseMock "lucassantoss1701/bank/internal/usecase/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
)

func GetBaseOriginAccount(t *testing.T) *entity.Account {
	originAccountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
	originAccountName := "lucas"
	originAccountCPF := "52849254088"
	originAccountsecret := "4578405"
	originAccountBalance := 100
	originAccountCreatedAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
	originAccount, err := entity.NewAccount(originAccountID, originAccountName, originAccountCPF, originAccountsecret, originAccountBalance, &originAccountCreatedAt)

	assert.Nil(t, err)
	assert.NotNil(t, originAccount)
	return originAccount
}

func GetBaseDestinationAccount(t *testing.T) *entity.Account {
	destionationAccountID := "d18551d3-cf13-49ec-b1dc-741a1f8715f6"
	destionationAccountName := "joao"
	destionationAccountCPF := "35768297090"
	destionationAccountSecret := "744637"
	destionationAccountBalance := 200
	destionationAccountCreatedAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
	destinationAccount, err := entity.NewAccount(destionationAccountID, destionationAccountName, destionationAccountCPF, destionationAccountSecret, destionationAccountBalance, &destionationAccountCreatedAt)

	assert.Nil(t, err)
	assert.NotNil(t, destinationAccount)
	return destinationAccount
}

func TestTransferHandler_Create(t *testing.T) {
	t.Run("Testing Create with success", func(t *testing.T) {
		transfer := mock.CreateTransfer()
		originAccount := GetBaseOriginAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		input := usecase.NewMakeTransferUseCaseInput(transfer.ID, originAccount.ID, destinationAccount.ID, transfer.Amount, transfer.CreatedAt)

		jsonData, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/transfers", bytes.NewBuffer(jsonData))

		req = req.WithContext(context.WithValue(req.Context(), web.AccountIDKey, originAccount.ID))

		recorder := httptest.NewRecorder()

		output := usecase.NewMakeTransferUseCaseOutput(&transfer)
		usecase := usecaseMock.NewMakeTransferUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebTransferHandler(usecase, nil)

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
	})

	t.Run("Testing Create when account id not exists in context", func(t *testing.T) {
		transfer := mock.CreateTransfer()
		originAccount := GetBaseOriginAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		input := usecase.NewMakeTransferUseCaseInput(transfer.ID, originAccount.ID, destinationAccount.ID, transfer.Amount, transfer.CreatedAt)

		jsonData, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/transfers", bytes.NewBuffer(jsonData))

		recorder := httptest.NewRecorder()

		output := usecase.NewMakeTransferUseCaseOutput(&transfer)
		usecase := usecaseMock.NewMakeTransferUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebTransferHandler(usecase, nil)

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Testing Create occurs error on decode body", func(t *testing.T) {
		transfer := mock.CreateTransfer()
		originAccount := GetBaseOriginAccount(t)

		req, _ := http.NewRequest("POST", "/transfers", bytes.NewBuffer([]byte("invalid json")))

		req = req.WithContext(context.WithValue(req.Context(), web.AccountIDKey, originAccount.ID))

		recorder := httptest.NewRecorder()

		output := usecase.NewMakeTransferUseCaseOutput(&transfer)
		usecase := usecaseMock.NewMakeTransferUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebTransferHandler(usecase, nil)

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Testing Create when usecase returns an error", func(t *testing.T) {
		transfer := mock.CreateTransfer()
		originAccount := GetBaseOriginAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		input := usecase.NewMakeTransferUseCaseInput(transfer.ID, originAccount.ID, destinationAccount.ID, transfer.Amount, transfer.CreatedAt)

		jsonData, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/transfers", bytes.NewBuffer(jsonData))

		req = req.WithContext(context.WithValue(req.Context(), web.AccountIDKey, originAccount.ID))

		recorder := httptest.NewRecorder()

		output := usecase.NewMakeTransferUseCaseOutput(&transfer)
		usecase := usecaseMock.NewMakeTransferUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, entity.NewErrorHandler(entity.INTERNAL_ERROR))

		handler := web.NewWebTransferHandler(usecase, nil)

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

}

func TestTransferHandler_FindByAccountID(t *testing.T) {
	t.Run("Testing FindByAccountID with success", func(t *testing.T) {
		transfer := mock.CreateTransfer()
		originAccount := GetBaseOriginAccount(t)

		req, _ := http.NewRequest("GET", "/transfers", nil)

		req = req.WithContext(context.WithValue(req.Context(), web.AccountIDKey, originAccount.ID))

		recorder := httptest.NewRecorder()

		var output []usecase.FindTransfersByAccountUseCaseOutput
		output = append(output, *usecase.NewFindTransfersByAccountUseCaseOutput(transfer))
		usecase := usecaseMock.NewFindTransfersByAccountUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebTransferHandler(nil, usecase)

		handler.FindByAccountID(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("Testing FindByAccountID  when account id not exists in context", func(t *testing.T) {
		transfer := mock.CreateTransfer()

		req, _ := http.NewRequest("GET", "/transfers", nil)

		recorder := httptest.NewRecorder()

		var output []usecase.FindTransfersByAccountUseCaseOutput
		output = append(output, *usecase.NewFindTransfersByAccountUseCaseOutput(transfer))
		usecase := usecaseMock.NewFindTransfersByAccountUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebTransferHandler(nil, usecase)

		handler.FindByAccountID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Testing FindByAccountID with with error(limit is wrong)", func(t *testing.T) {
		transfer := mock.CreateTransfer()
		originAccount := GetBaseOriginAccount(t)

		req, _ := http.NewRequest("GET", "/transfers?limit=abc&offset=0", nil)

		req = req.WithContext(context.WithValue(req.Context(), web.AccountIDKey, originAccount.ID))

		recorder := httptest.NewRecorder()

		var output []usecase.FindTransfersByAccountUseCaseOutput
		output = append(output, *usecase.NewFindTransfersByAccountUseCaseOutput(transfer))
		usecase := usecaseMock.NewFindTransfersByAccountUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebTransferHandler(nil, usecase)

		handler.FindByAccountID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Testing FindByAccountID with with error(offset is wrong)", func(t *testing.T) {
		transfer := mock.CreateTransfer()
		originAccount := GetBaseOriginAccount(t)

		req, _ := http.NewRequest("GET", "/transfers?limit=10&offset=dfs", nil)

		req = req.WithContext(context.WithValue(req.Context(), web.AccountIDKey, originAccount.ID))

		recorder := httptest.NewRecorder()

		var output []usecase.FindTransfersByAccountUseCaseOutput
		output = append(output, *usecase.NewFindTransfersByAccountUseCaseOutput(transfer))
		usecase := usecaseMock.NewFindTransfersByAccountUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebTransferHandler(nil, usecase)

		handler.FindByAccountID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Testing FindByAccountID when usecase returns an error", func(t *testing.T) {
		transfer := mock.CreateTransfer()
		originAccount := GetBaseOriginAccount(t)

		req, _ := http.NewRequest("GET", "/transfers", nil)

		req = req.WithContext(context.WithValue(req.Context(), web.AccountIDKey, originAccount.ID))

		recorder := httptest.NewRecorder()

		var output []usecase.FindTransfersByAccountUseCaseOutput
		output = append(output, *usecase.NewFindTransfersByAccountUseCaseOutput(transfer))
		usecase := usecaseMock.NewFindTransfersByAccountUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, entity.NewErrorHandler(entity.INTERNAL_ERROR))

		handler := web.NewWebTransferHandler(nil, usecase)

		handler.FindByAccountID(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

}
