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

// @Summary     Create
// @Description Create account
// @Tags        accounts
// @Produce     json
// @Param       body body usecase.CreateAccountUseCaseInput true "create account request vody"
// @Success     201 {object} usecase.CreateAccountUseCaseOutput
// @Failure     400,401,404,500,422
// @Router /accounts [post]
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

// @Summary     Find
// @Description Find accounts by param
// @Tags        accounts
// @Produce     json
// @Param       limit query int false "number of items to be returned per page"
// @Param       offset query int false "page offset"
// @Success     200 {array} usecase.FindAccountUseCaseOutput
// @Failure     400,401,404,500
// @Security    ApiKeyAuth
// @Router /accounts [get]
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

// @Summary     Find
// @Description Find balance of a specific accounts
// @Tags        accounts
// @Produce     json
// @Success     200 {array} usecase.FindBalanceByAccountUseCaseOutput
// @Param       account_id path string true "account_id"
// @Failure     400,401,404,500
// @Security    ApiKeyAuth
// @Router /accounts/{account_id}/balance [get]
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

// @Summary     Login
// @Description Login checks that the user can use the API and returns a token
// @Tags        accounts
// @Produce     json
// @Param       body body usecase.LoginUseCaseInput true "login request body"
// @Success     200 {object} usecase.LoginUseCaseOutput
// @Failure     400,401,404,500
// @Router /login [post]
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
