package web

import (
	"encoding/json"
	"errors"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/infra/web/responses"
	"lucassantoss1701/bank/internal/usecase"
	"net/http"
	"strconv"
	"time"
)

type ContextKey string

const AccountIDKey ContextKey = "account_id"

type WebTransferHandler struct {
	accountRepository  entity.AccountRepository
	transferRepostiory entity.TransferRepository
	repository         entity.Repository
}

func NewWebTransferHandler(accountRepository entity.AccountRepository, transferRepostiory entity.TransferRepository, repository entity.Repository) *WebTransferHandler {
	return &WebTransferHandler{
		accountRepository:  accountRepository,
		transferRepostiory: transferRepostiory,
		repository:         repository,
	}
}

func (h *WebTransferHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accountID, ok := ctx.Value(AccountIDKey).(string)
	if !ok {
		responses.Err(w, errors.New("account_id not found in context"))
		return
	}

	var dto usecase.MakeTransferUseCaseInput
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		responses.Err(w, err)
		return
	}

	createdAt := time.Now()
	input := usecase.NewMakeTransferUseCaseInput(dto.ID, accountID, dto.DestinationAccount.ID, dto.Amount, &createdAt)

	makeTransfer := usecase.NewMakeTransferUseCase(h.accountRepository, h.transferRepostiory, h.repository)
	output, err := makeTransfer.Execute(ctx, input)
	if err != nil {
		responses.Err(w, err)
		return
	}

	responses.Success(w, http.StatusCreated, output)
}

func (h *WebTransferHandler) FindByAccountID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	accountID, ok := ctx.Value(AccountIDKey).(string)
	if !ok {
		responses.Err(w, errors.New("account_id not found in context"))
		return
	}

	limit := 0
	offSet := 0

	queryParams := r.URL.Query()

	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			responses.Err(w, err)
			return
		}
	}

	if limit == 0 {
		limit = 20
	}

	offSetStr := queryParams.Get("offset")
	if offSetStr != "" {
		offSet, err = strconv.Atoi(offSetStr)
		if err != nil {
			responses.Err(w, err)
			return
		}
	}

	input := usecase.NewFindTransfersByAccountUseCaseInput(accountID, limit, offSet)

	findTransfersByAccountUseCase := usecase.NewFindTransfersByAccountUseCase(h.transferRepostiory)
	output, err := findTransfersByAccountUseCase.Execute(ctx, input)
	if err != nil {
		responses.Err(w, err)
		return
	}

	responses.Success(w, http.StatusOK, output)
}
