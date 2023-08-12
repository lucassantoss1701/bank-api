package web

import (
	"encoding/json"
	"lucassantoss1701/bank/configs"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/infra/web/responses"
	"lucassantoss1701/bank/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type WebAccountHandler struct {
	createAccount usecase.ICreateAccountUseCase
	findAccount   usecase.IFindAccountUseCase
	findBalance   usecase.IFindBalanceByAccountUseCase
	login         usecase.ILoginUseCase
}

func NewWebAccountHandler(createAccount usecase.ICreateAccountUseCase, findAccount usecase.IFindAccountUseCase, findBalance usecase.IFindBalanceByAccountUseCase, login usecase.ILoginUseCase) *WebAccountHandler {
	return &WebAccountHandler{
		createAccount: createAccount,
		findAccount:   findAccount,
		findBalance:   findBalance,
		login:         login,
	}
}

func (h *WebAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto usecase.CreateAccountUseCaseInput
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		responses.Err(w, entity.NewErrorHandler(entity.BAD_REQUEST).Add(err.Error()))
		return
	}

	createdAt := time.Now()
	dto.CreatedAt = &createdAt

	output, err := h.createAccount.Execute(ctx, &dto)
	if err != nil {
		responses.Err(w, err)
		return
	}

	responses.Success(w, http.StatusCreated, output)
}

func (h *WebAccountHandler) Find(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	limit := 0
	offSet := 0

	queryParams := r.URL.Query()

	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			responses.Err(w, entity.NewErrorHandler(entity.BAD_REQUEST).Add(err.Error()))
			return
		}
		if limit == 0 {
			limit = 20
		}
	}

	offSetStr := queryParams.Get("offset")
	if offSetStr != "" {
		offSet, err = strconv.Atoi(offSetStr)
		if err != nil {
			responses.Err(w, entity.NewErrorHandler(entity.BAD_REQUEST).Add(err.Error()))
			return
		}
	}

	input := usecase.NewFindAccountUseCaseInput(limit, offSet)
	output, err := h.findAccount.Execute(ctx, input)
	if err != nil {
		responses.Err(w, err)
		return
	}
	responses.Success(w, http.StatusOK, output)
}

func (h *WebAccountHandler) FindBalanceByAccount(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	accountID := chi.URLParam(r, "account_id")

	input := usecase.NewFindBalanceByAccountUseCaseInput(accountID)
	output, err := h.findBalance.Execute(ctx, input)
	if err != nil {
		responses.Err(w, err)
		return
	}

	responses.Success(w, http.StatusOK, output)
}

func (h *WebAccountHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto usecase.LoginUseCaseInput
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		responses.Err(w, entity.NewErrorHandler(entity.BAD_REQUEST).Add(err.Error()))
		return
	}

	secretJWT := configs.Get().Security.Secret
	input := usecase.NewLoginUseCaseInput(dto.CPF, dto.Secret, secretJWT)

	output, err := h.login.Execute(ctx, input)
	if err != nil {
		responses.Err(w, err)
		return
	}

	responses.Success(w, http.StatusOK, output)

}
