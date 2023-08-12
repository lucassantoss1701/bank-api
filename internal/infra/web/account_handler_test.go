package web_test

import (
	"bytes"
	"encoding/json"
	"lucassantoss1701/bank/configs"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/entity/mock"
	"lucassantoss1701/bank/internal/infra/web"
	"lucassantoss1701/bank/internal/usecase"
	usecaseMock "lucassantoss1701/bank/internal/usecase/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
)

func init() {
	configs.Load()
}
func TestAccountHandler_Create(t *testing.T) {
	t.Run("Testing create with success", func(t *testing.T) {

		account := mock.CreateAccount()

		input := usecase.NewCreateAccountUseCaseInput(account.ID, account.Name, account.CPF, account.Secret, account.Balance, *account.CreatedAt)

		jsonData, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonData))

		recorder := httptest.NewRecorder()

		output := usecase.NewCreateAccountUseCaseOutput(account.ID, account.Name, account.Balance, account.CreatedAt)
		usecase := usecaseMock.NewCreateAccountUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebAccountHandler(usecase, nil, nil, nil)

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
	})

	t.Run("Testing create when occurs error on decode body", func(t *testing.T) {

		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte("invalid json")))

		recorder := httptest.NewRecorder()

		usecase := usecaseMock.NewCreateAccountUseCaseMock()
		handler := web.NewWebAccountHandler(usecase, nil, nil, nil)

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Testing create when execute returns a error", func(t *testing.T) {

		account := mock.CreateAccount()

		input := usecase.NewCreateAccountUseCaseInput(account.ID, account.Name, account.CPF, account.Secret, account.Balance, *account.CreatedAt)

		jsonData, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonData))

		recorder := httptest.NewRecorder()

		output := usecase.NewCreateAccountUseCaseOutput(account.ID, account.Name, account.Balance, account.CreatedAt)
		usecase := usecaseMock.NewCreateAccountUseCaseMock()

		usecase.On("Execute", req.Context(), testify.Anything).Return(output, entity.NewErrorHandler(entity.INTERNAL_ERROR))

		handler := web.NewWebAccountHandler(usecase, nil, nil, nil)

		handler.Create(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestAccountHandler_Find(t *testing.T) {
	t.Run("Testing find with success", func(t *testing.T) {
		account := mock.CreateAccount()
		req, _ := http.NewRequest("GET", "/accounts?limit=0&offset=0", nil)
		recorder := httptest.NewRecorder()
		var output []usecase.FindAccountUseCaseOutput
		output = append(output, *usecase.NewFindAccountUseCaseOutput(account.ID, account.Name, account.Balance, account.CreatedAt))
		usecase := usecaseMock.NewFindAccountUseCaseMock()

		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebAccountHandler(nil, usecase, nil, nil)

		handler.Find(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("Testing find with error(limit is wrong)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/accounts?limit=abc&offset=0", nil)
		recorder := httptest.NewRecorder()
		usecase := usecaseMock.NewFindAccountUseCaseMock()
		handler := web.NewWebAccountHandler(nil, usecase, nil, nil)
		handler.Find(recorder, req)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Testing find with error(offset is wrong)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/accounts?limit=10&offset=dfs", nil)
		recorder := httptest.NewRecorder()
		usecase := usecaseMock.NewFindAccountUseCaseMock()
		handler := web.NewWebAccountHandler(nil, usecase, nil, nil)
		handler.Find(recorder, req)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Testing find when usecase returns an error", func(t *testing.T) {
		account := mock.CreateAccount()
		req, _ := http.NewRequest("GET", "/accounts?limit=0&offset=0", nil)
		recorder := httptest.NewRecorder()
		var output []usecase.FindAccountUseCaseOutput
		output = append(output, *usecase.NewFindAccountUseCaseOutput(account.ID, account.Name, account.Balance, account.CreatedAt))
		usecase := usecaseMock.NewFindAccountUseCaseMock()

		usecase.On("Execute", req.Context(), testify.Anything).Return(output, entity.NewErrorHandler(entity.INTERNAL_ERROR))

		handler := web.NewWebAccountHandler(nil, usecase, nil, nil)

		handler.Find(recorder, req)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

}

func TestAccountHandler_FindBalanceByAccount(t *testing.T) {
	t.Run("Testing FindBalanceByAccount with success", func(t *testing.T) {
		account := mock.CreateAccount()
		req, _ := http.NewRequest("GET", "/accounts/id", nil)
		recorder := httptest.NewRecorder()
		output := usecase.NewFindBalanceByAccountUseCaseOutput(account.Balance)
		usecase := usecaseMock.NewFindBalanceByAccountUseCaseMock()

		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebAccountHandler(nil, nil, usecase, nil)

		handler.FindBalanceByAccount(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("Testing FindBalanceByAccount  when usecase returns an error", func(t *testing.T) {
		account := mock.CreateAccount()
		req, _ := http.NewRequest("GET", "/accounts/id", nil)
		recorder := httptest.NewRecorder()
		output := usecase.NewFindBalanceByAccountUseCaseOutput(account.Balance)
		usecase := usecaseMock.NewFindBalanceByAccountUseCaseMock()

		usecase.On("Execute", req.Context(), testify.Anything).Return(output, entity.NewErrorHandler(entity.INTERNAL_ERROR))

		handler := web.NewWebAccountHandler(nil, nil, usecase, nil)

		handler.FindBalanceByAccount(recorder, req)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestAccountHandler_Login(t *testing.T) {
	t.Run("Testing Login with success", func(t *testing.T) {

		account := mock.CreateAccount()

		input := usecase.NewLoginUseCaseInput(account.CPF, account.Secret, "")

		jsonData, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))

		recorder := httptest.NewRecorder()

		output := usecase.NewLoginUseCaseOutput(&account, "")

		usecase := usecaseMock.NewLoginUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, nil)

		handler := web.NewWebAccountHandler(nil, nil, nil, usecase)

		handler.Login(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("Testing Login when occurs error on decode body", func(t *testing.T) {

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))

		recorder := httptest.NewRecorder()

		usecase := usecaseMock.NewLoginUseCaseMock()
		handler := web.NewWebAccountHandler(nil, nil, nil, usecase)

		handler.Login(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Testing Login when execute returns a error", func(t *testing.T) {

		account := mock.CreateAccount()

		input := usecase.NewLoginUseCaseInput(account.CPF, account.Secret, "")

		jsonData, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))

		recorder := httptest.NewRecorder()

		output := usecase.NewLoginUseCaseOutput(&account, "")
		usecase := usecaseMock.NewLoginUseCaseMock()
		usecase.On("Execute", req.Context(), testify.Anything).Return(output, entity.NewErrorHandler(entity.INTERNAL_ERROR))

		handler := web.NewWebAccountHandler(nil, nil, nil, usecase)

		handler.Login(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}
