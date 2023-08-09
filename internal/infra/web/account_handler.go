package web

import (
	"encoding/json"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/infra/web/responses"
	"lucassantoss1701/bank/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type WebAccountHandler struct {
	accountRepostiory entity.AccountRepository
}

func NewWebAccountHandler(accountRepostiory entity.AccountRepository) *WebAccountHandler {
	return &WebAccountHandler{
		accountRepostiory: accountRepostiory,
	}
}

func (h *WebAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto usecase.CreateAccountUseCaseInput
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		responses.Err(w, err)
		return
	}

	createdAt := time.Now()
	dto.CreatedAt = &createdAt

	createAccount := usecase.NewCreateAccountUseCase(h.accountRepostiory)
	output, err := createAccount.Execute(ctx, &dto)
	if err != nil {
		responses.Err(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		responses.Err(w, err)
		return
	}

}

func (h *WebAccountHandler) Find(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	limit := 0
	offSet := 0

	findAccounts := usecase.NewFindAccountUseCase(h.accountRepostiory)
	queryParams := r.URL.Query()

	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			responses.Err(w, err)
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
			responses.Err(w, err)
			return
		}
	}

	input := usecase.NewFindAccountUseCaseInput(limit, offSet)
	output, err := findAccounts.Execute(ctx, input)
	if err != nil {
		responses.Err(w, err)
		return
	}
	responses.JSON(w, http.StatusOK, output)
}

func (h *WebAccountHandler) FindBalanceByAccount(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	accountID := chi.URLParam(r, "account_id")
	findBalance := usecase.NewFindBalanceByAccountUseCase(h.accountRepostiory)

	input := usecase.NewFindBalanceByAccountUseCaseInput(accountID)
	output, err := findBalance.Execute(ctx, input)
	if err != nil {
		responses.Err(w, err)
		return
	}

	responses.JSON(w, http.StatusOK, output)
}
