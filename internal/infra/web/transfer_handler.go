package web

import (
	"encoding/json"
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
	makeTransfer          usecase.IMakeTransferUseCase
	findTransferByAccount usecase.IFindTransfersByAccountUseCase
}

func NewWebTransferHandler(makeTransfer usecase.IMakeTransferUseCase, findTransferByAccount usecase.IFindTransfersByAccountUseCase) *WebTransferHandler {
	return &WebTransferHandler{
		makeTransfer:          makeTransfer,
		findTransferByAccount: findTransferByAccount,
	}
}

// @Summary     Create transfer
// @Description Create transfer between two accounts
// @Tags        transfers
// @Produce     json
// @Param       body body usecase.MakeTransferUseCaseInput true "make transfer request body"
// @Success     201 {object} usecase.MakeTransferUseCaseOutput
// @Failure     400,401,404,500,422
// @Security    ApiKeyAuth
// @Router /transfers [post]
func (h *WebTransferHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accountID, ok := ctx.Value(AccountIDKey).(string)
	if !ok {
		responses.Err(w, entity.NewErrorHandler(entity.BAD_REQUEST).Add("account_id not found in context"))
		return
	}

	var dto usecase.MakeTransferUseCaseInput
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		responses.Err(w, entity.NewErrorHandler(entity.BAD_REQUEST).Add(err.Error()))
		return
	}

	createdAt := time.Now()
	input := usecase.NewMakeTransferUseCaseInput(dto.ID, accountID, dto.DestinationAccount.ID, dto.Amount, &createdAt)

	output, err := h.makeTransfer.Execute(ctx, input)
	if err != nil {
		responses.Err(w, err)
		return
	}

	responses.Success(w, http.StatusCreated, output)
}

// @Summary     Find transfers by account
// @Description Find transfers from an account(user needs to be authenticated)
// @Tags        transfers
// @Produce     json
// @Param       limit query int false "number of items to be returned per page"
// @Param       offset query int false "page offset"
// @Success     200 {array} usecase.MakeTransferUseCaseOutput
// @Failure     400,401,404,500
// @Security    ApiKeyAuth
// @Router /transfers [get]
func (h *WebTransferHandler) FindByAccountID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	accountID, ok := ctx.Value(AccountIDKey).(string)
	if !ok {
		responses.Err(w, entity.NewErrorHandler(entity.BAD_REQUEST).Add("account_id not found in context"))
		return
	}

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
	}

	if limit == 0 {
		limit = 20
	}

	offSetStr := queryParams.Get("offset")
	if offSetStr != "" {
		offSet, err = strconv.Atoi(offSetStr)
		if err != nil {
			responses.Err(w, entity.NewErrorHandler(entity.BAD_REQUEST).Add(err.Error()))
			return
		}
	}

	input := usecase.NewFindTransfersByAccountUseCaseInput(accountID, limit, offSet)
	output, err := h.findTransferByAccount.Execute(ctx, input)
	if err != nil {
		responses.Err(w, err)
		return
	}

	responses.Success(w, http.StatusOK, output)
}
