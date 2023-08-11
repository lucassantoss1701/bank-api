package usecase_test

import (
	"context"
	"errors"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/entity/mock"
	"lucassantoss1701/bank/internal/usecase"
	"testing"
	"time"

	testify "github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
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

func ReturnTransferPointer(ctx context.Context, transfer *entity.Transfer, transactionHandlers ...entity.TransactionHandler) *entity.Transfer {
	return transfer
}

func TestMakeTransferUseCase_Execute(t *testing.T) {
	t.Run("Testing MakeTransferUseCase when have success on make transfer", func(t *testing.T) {
		ctx := context.Background()
		amount := 50

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		originAccountAfterTransfer := *originAccount
		destinationAccountAfterTransfer := *destinationAccount

		originAccountAfterTransfer.Balance -= amount      // NewBalance = 50
		destinationAccountAfterTransfer.Balance += amount // NewBalance  = 250

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(*originAccount, nil)
		accountRepository.On("FindByID", ctx, destinationAccount.ID).Return(*destinationAccount, nil)
		accountRepository.On("UpdateBalance", ctx, originAccount.ID, originAccountAfterTransfer.Balance, testify.Anything).Return(originAccountAfterTransfer, nil)
		accountRepository.On("UpdateBalance", ctx, destinationAccount.ID, destinationAccountAfterTransfer.Balance, testify.Anything).Return(destinationAccountAfterTransfer, nil)

		transactionHandler := mock.NewTransactionHandlerMock()

		repository := mock.NewRepositoryMock()
		repository.On("BeginTx", ctx).Return(transactionHandler, nil)
		repository.On("CommitTx", transactionHandler).Return(nil)

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"
		createdAt := time.Date(2023, 8, 7, 10, 00, 00, 00, time.UTC)
		transferRepository := mock.NewTransferRepositoryMock()
		transfer, err := entity.NewTransfer(transferID, originAccount, destinationAccount, amount, &createdAt)
		assert.Nil(t, err)

		transferAfterTransaction := *transfer
		transferAfterTransaction.OriginAccount = &originAccountAfterTransfer
		transferAfterTransaction.DestinationAccount = &destinationAccountAfterTransfer
		transferRepository.On("Create", ctx, &transferAfterTransaction, testify.Anything).Return(transferAfterTransaction, nil)

		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, &createdAt)
		output, err := makeTransferUseCase.Execute(ctx, input)

		assert.Nil(t, err)
		assert.NotNil(t, output)

		assert.Equal(t, transferID, output.ID)
		assert.Equal(t, amount, output.Amount)
		assert.Equal(t, createdAt, *output.CreatedAt)

		assert.Equal(t, originAccount.ID, output.OriginAccount.ID)
		assert.Equal(t, originAccount.Name, output.OriginAccount.Name)

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertCalled(t, "FindByID", ctx, destinationAccount.ID)
		accountRepository.AssertCalled(t, "UpdateBalance", ctx, originAccount.ID, 50, testify.Anything)
		accountRepository.AssertCalled(t, "UpdateBalance", ctx, destinationAccount.ID, 250, testify.Anything)
		transferRepository.AssertCalled(t, "Create", ctx, &transferAfterTransaction, testify.Anything)
		repository.AssertCalled(t, "BeginTx", ctx)
		repository.AssertCalled(t, "CommitTx", testify.Anything)
		repository.AssertNotCalled(t, "RollbackTx", testify.Anything)
	})

	t.Run("Testing MakeTransferUseCase when origin account not found", func(t *testing.T) {
		ctx := context.Background()
		amount := 50

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(entity.Account{}, errors.New("origin account not found"))

		transferRepository := mock.NewTransferRepositoryMock()
		repository := mock.NewRepositoryMock()

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"
		createdAt := time.Date(2023, 8, 7, 10, 00, 00, 00, time.UTC)
		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, &createdAt)
		output, err := makeTransferUseCase.Execute(ctx, input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "origin account not found", err.Error())

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, originAccount.ID, testify.Anything)
		accountRepository.AssertNotCalled(t, "FindByID", ctx, destinationAccount.ID)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, destinationAccount.ID, testify.Anything)
		transferRepository.AssertNotCalled(t, "Create", ctx, testify.Anything, testify.Anything)
		repository.AssertNotCalled(t, "BeginTx", ctx)
		repository.AssertNotCalled(t, "CommitTx", testify.Anything)
		repository.AssertNotCalled(t, "RollbackTx", testify.Anything)
	})

	t.Run("Testing MakeTransferUseCase when destination account not found", func(t *testing.T) {
		ctx := context.Background()
		amount := 50

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(*originAccount, nil)
		accountRepository.On("FindByID", ctx, destinationAccount.ID).Return(entity.Account{}, errors.New("destination account not found"))

		transferRepository := mock.NewTransferRepositoryMock()
		repository := mock.NewRepositoryMock()

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"
		createdAt := time.Date(2023, 8, 7, 10, 00, 00, 00, time.UTC)
		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, &createdAt)
		output, err := makeTransferUseCase.Execute(ctx, input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "destination account not found", err.Error())

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertCalled(t, "FindByID", ctx, destinationAccount.ID)
		transferRepository.AssertNotCalled(t, "Create", ctx, testify.Anything, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, originAccount.ID, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, destinationAccount.ID, testify.Anything)
		repository.AssertNotCalled(t, "BeginTx", ctx)
		repository.AssertNotCalled(t, "CommitTx", testify.Anything)
		repository.AssertNotCalled(t, "RollbackTx", testify.Anything)

	})

	t.Run("Testing MakeTransferUseCase when transfer amount exceeds balance", func(t *testing.T) {
		ctx := context.Background()
		amount := 150

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(*originAccount, nil)
		accountRepository.On("FindByID", ctx, destinationAccount.ID).Return(*destinationAccount, nil)

		transferRepository := mock.NewTransferRepositoryMock()
		repository := mock.NewRepositoryMock()

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"
		createdAt := time.Date(2023, 8, 7, 10, 00, 00, 00, time.UTC)
		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, &createdAt)
		output, err := makeTransferUseCase.Execute(ctx, input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "error on update balance of origin account: new balance cannot be minor than 0", err.Error())

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertCalled(t, "FindByID", ctx, destinationAccount.ID)
		transferRepository.AssertNotCalled(t, "Create", ctx, testify.Anything, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, originAccount.ID, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, destinationAccount.ID, testify.Anything)
		repository.AssertNotCalled(t, "BeginTx", ctx)
		repository.AssertNotCalled(t, "CommitTx", testify.Anything)
		repository.AssertNotCalled(t, "RollbackTx", testify.Anything)
	})

	t.Run("Testing MakeTransferUseCase when newTransfer returns an error", func(t *testing.T) {
		ctx := context.Background()
		amount := 50

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(*originAccount, nil)
		accountRepository.On("FindByID", ctx, destinationAccount.ID).Return(*destinationAccount, nil)

		transferRepository := mock.NewTransferRepositoryMock()
		repository := mock.NewRepositoryMock()

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"

		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, nil)
		output, err := makeTransferUseCase.Execute(ctx, input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "created at cannot be nil", err.Error())

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertCalled(t, "FindByID", ctx, destinationAccount.ID)
		transferRepository.AssertNotCalled(t, "Create", ctx, testify.Anything, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, originAccount.ID, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, destinationAccount.ID, testify.Anything)
		repository.AssertNotCalled(t, "BeginTx", ctx)
		repository.AssertNotCalled(t, "CommitTx", testify.Anything)
		repository.AssertNotCalled(t, "RollbackTx", testify.Anything)
	})

	t.Run("Testing MakeTransferUseCase when beginTx returns an error", func(t *testing.T) {
		ctx := context.Background()
		amount := 50

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(*originAccount, nil)
		accountRepository.On("FindByID", ctx, destinationAccount.ID).Return(*destinationAccount, nil)

		transferRepository := mock.NewTransferRepositoryMock()

		repository := mock.NewRepositoryMock()
		transactionHandler := mock.NewTransactionHandlerMock()
		repository.On("BeginTx", ctx).Return(transactionHandler, errors.New("error on begin transaction"))

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"

		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		createdAt := time.Date(2023, 8, 7, 10, 00, 00, 00, time.UTC)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, &createdAt)
		output, err := makeTransferUseCase.Execute(ctx, input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "error on begin transaction", err.Error())

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertCalled(t, "FindByID", ctx, destinationAccount.ID)
		transferRepository.AssertNotCalled(t, "Create", ctx, testify.Anything, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, originAccount.ID, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, destinationAccount.ID, testify.Anything)
		repository.AssertCalled(t, "BeginTx", ctx)
		repository.AssertNotCalled(t, "CommitTx", testify.Anything)
		repository.AssertNotCalled(t, "RollbackTx", testify.Anything)
	})

	t.Run("Testing MakeTransferUseCase when beginTx returns an error", func(t *testing.T) {
		ctx := context.Background()
		amount := 50

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(*originAccount, nil)
		accountRepository.On("FindByID", ctx, destinationAccount.ID).Return(*destinationAccount, nil)

		transferRepository := mock.NewTransferRepositoryMock()

		repository := mock.NewRepositoryMock()
		transactionHandler := mock.NewTransactionHandlerMock()
		repository.On("BeginTx", ctx).Return(transactionHandler, errors.New("error on begin transaction"))

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"

		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		createdAt := time.Date(2023, 8, 7, 10, 00, 00, 00, time.UTC)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, &createdAt)
		output, err := makeTransferUseCase.Execute(ctx, input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "error on begin transaction", err.Error())

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertCalled(t, "FindByID", ctx, destinationAccount.ID)
		transferRepository.AssertNotCalled(t, "Create", ctx, testify.Anything, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, originAccount.ID, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, destinationAccount.ID, testify.Anything)
		repository.AssertCalled(t, "BeginTx", ctx)
		repository.AssertNotCalled(t, "CommitTx", testify.Anything)
		repository.AssertNotCalled(t, "RollbackTx", testify.Anything)

	})

	t.Run("Testing MakeTransferUseCase when create transfer returns an error", func(t *testing.T) {
		ctx := context.Background()
		amount := 50

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		originAccountAfterTransfer := *originAccount
		destinationAccountAfterTransfer := *destinationAccount

		originAccountAfterTransfer.Balance -= amount      // NewBalance = 50
		destinationAccountAfterTransfer.Balance += amount // NewBalance  = 250

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(*originAccount, nil)
		accountRepository.On("FindByID", ctx, destinationAccount.ID).Return(*destinationAccount, nil)
		accountRepository.On("UpdateBalance", ctx, originAccount.ID, originAccountAfterTransfer.Balance, testify.Anything).Return(originAccountAfterTransfer, nil)
		accountRepository.On("UpdateBalance", ctx, destinationAccount.ID, destinationAccountAfterTransfer.Balance, testify.Anything).Return(destinationAccountAfterTransfer, nil)

		transactionHandler := mock.NewTransactionHandlerMock()

		repository := mock.NewRepositoryMock()
		repository.On("BeginTx", ctx).Return(transactionHandler, nil)
		repository.On("CommitTx", transactionHandler).Return(nil)
		repository.On("RollbackTx", transactionHandler).Return(nil)

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"
		createdAt := time.Date(2023, 8, 7, 10, 00, 00, 00, time.UTC)
		transferRepository := mock.NewTransferRepositoryMock()
		transfer, err := entity.NewTransfer(transferID, originAccount, destinationAccount, amount, &createdAt)
		assert.Nil(t, err)

		transferAfterTransaction := *transfer
		transferAfterTransaction.OriginAccount = &originAccountAfterTransfer
		transferAfterTransaction.DestinationAccount = &destinationAccountAfterTransfer

		var returnedTransaction entity.Transfer

		transferRepository.On("Create", ctx, &transferAfterTransaction, testify.Anything).Return(returnedTransaction, errors.New("error on create transfer"))

		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, &createdAt)
		output, err := makeTransferUseCase.Execute(ctx, input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "error on create transfer", err.Error())

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertCalled(t, "FindByID", ctx, destinationAccount.ID)
		transferRepository.AssertCalled(t, "Create", ctx, testify.Anything, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, originAccount.ID, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, destinationAccount.ID, testify.Anything)
		repository.AssertCalled(t, "BeginTx", ctx)
		repository.AssertNotCalled(t, "CommitTx", testify.Anything)
		repository.AssertCalled(t, "RollbackTx", testify.Anything)
	})

	t.Run("Testing MakeTransferUseCase when update origin account balance returns a error", func(t *testing.T) {
		ctx := context.Background()
		amount := 50

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		originAccountAfterTransfer := *originAccount
		destinationAccountAfterTransfer := *destinationAccount

		originAccountAfterTransfer.Balance -= amount      // NewBalance = 50
		destinationAccountAfterTransfer.Balance += amount // NewBalance  = 250

		var returnedOriginAccount entity.Account

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(*originAccount, nil)
		accountRepository.On("FindByID", ctx, destinationAccount.ID).Return(*destinationAccount, nil)
		accountRepository.On("UpdateBalance", ctx, originAccount.ID, originAccountAfterTransfer.Balance, testify.Anything).Return(returnedOriginAccount, errors.New("error on update origin account balance"))
		accountRepository.On("UpdateBalance", ctx, destinationAccount.ID, destinationAccountAfterTransfer.Balance, testify.Anything).Return(destinationAccountAfterTransfer, nil)

		transactionHandler := mock.NewTransactionHandlerMock()

		repository := mock.NewRepositoryMock()
		repository.On("BeginTx", ctx).Return(transactionHandler, nil)
		repository.On("CommitTx", transactionHandler).Return(nil)
		repository.On("RollbackTx", transactionHandler).Return(nil)

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"
		createdAt := time.Date(2023, 8, 7, 10, 00, 00, 00, time.UTC)
		transferRepository := mock.NewTransferRepositoryMock()
		transfer, err := entity.NewTransfer(transferID, originAccount, destinationAccount, amount, &createdAt)
		assert.Nil(t, err)

		transferAfterTransaction := *transfer
		transferAfterTransaction.OriginAccount = &originAccountAfterTransfer
		transferAfterTransaction.DestinationAccount = &destinationAccountAfterTransfer

		transferRepository.On("Create", ctx, &transferAfterTransaction, testify.Anything).Return(transferAfterTransaction, nil)

		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, &createdAt)
		output, err := makeTransferUseCase.Execute(ctx, input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "error on update origin account balance", err.Error())

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertCalled(t, "FindByID", ctx, destinationAccount.ID)
		transferRepository.AssertCalled(t, "Create", ctx, testify.Anything, testify.Anything)
		accountRepository.AssertCalled(t, "UpdateBalance", ctx, originAccount.ID, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, destinationAccount.ID, testify.Anything)
		repository.AssertCalled(t, "BeginTx", ctx)
		repository.AssertNotCalled(t, "CommitTx", testify.Anything)
		repository.AssertCalled(t, "RollbackTx", testify.Anything)
	})

	t.Run("Testing MakeTransferUseCase when update destination account balance returns a error", func(t *testing.T) {
		ctx := context.Background()
		amount := 50

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		originAccountAfterTransfer := *originAccount
		destinationAccountAfterTransfer := *destinationAccount

		originAccountAfterTransfer.Balance -= amount      // NewBalance = 50
		destinationAccountAfterTransfer.Balance += amount // NewBalance  = 250

		var returnedDestinationAccount entity.Account

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(*originAccount, nil)
		accountRepository.On("FindByID", ctx, destinationAccount.ID).Return(*destinationAccount, nil)
		accountRepository.On("UpdateBalance", ctx, originAccount.ID, originAccountAfterTransfer.Balance, testify.Anything).Return(originAccountAfterTransfer, nil)
		accountRepository.On("UpdateBalance", ctx, destinationAccount.ID, destinationAccountAfterTransfer.Balance, testify.Anything).Return(returnedDestinationAccount, errors.New("error on update destination account balance"))

		transactionHandler := mock.NewTransactionHandlerMock()

		repository := mock.NewRepositoryMock()
		repository.On("BeginTx", ctx).Return(transactionHandler, nil)
		repository.On("CommitTx", transactionHandler).Return(nil)
		repository.On("RollbackTx", transactionHandler).Return(nil)

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"
		createdAt := time.Date(2023, 8, 7, 10, 00, 00, 00, time.UTC)
		transferRepository := mock.NewTransferRepositoryMock()
		transfer, err := entity.NewTransfer(transferID, originAccount, destinationAccount, amount, &createdAt)
		assert.Nil(t, err)

		transferAfterTransaction := *transfer
		transferAfterTransaction.OriginAccount = &originAccountAfterTransfer
		transferAfterTransaction.DestinationAccount = &destinationAccountAfterTransfer

		transferRepository.On("Create", ctx, &transferAfterTransaction, testify.Anything).Return(transferAfterTransaction, nil)

		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, &createdAt)
		output, err := makeTransferUseCase.Execute(ctx, input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "error on update destination account balance", err.Error())

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertCalled(t, "FindByID", ctx, destinationAccount.ID)
		transferRepository.AssertCalled(t, "Create", ctx, testify.Anything, testify.Anything)
		accountRepository.AssertCalled(t, "UpdateBalance", ctx, originAccount.ID, testify.Anything)
		accountRepository.AssertCalled(t, "UpdateBalance", ctx, destinationAccount.ID, testify.Anything)
		repository.AssertCalled(t, "BeginTx", ctx)
		repository.AssertNotCalled(t, "CommitTx", testify.Anything)
		repository.AssertCalled(t, "RollbackTx", testify.Anything)
	})

	t.Run("Testing MakeTransferUseCase when update destination account balance returns a error", func(t *testing.T) {
		ctx := context.Background()
		amount := 50

		originAccount := GetBaseOriginAccount(t)           // Balance = 100
		destinationAccount := GetBaseDestinationAccount(t) // Balance = 200

		originAccountAfterTransfer := *originAccount
		destinationAccountAfterTransfer := *destinationAccount

		originAccountAfterTransfer.Balance -= amount      // NewBalance = 50
		destinationAccountAfterTransfer.Balance += amount // NewBalance  = 250

		accountRepository := mock.NewAccountRepositoryMock()
		accountRepository.On("FindByID", ctx, originAccount.ID).Return(*originAccount, nil)
		accountRepository.On("FindByID", ctx, destinationAccount.ID).Return(*destinationAccount, nil)
		accountRepository.On("UpdateBalance", ctx, originAccount.ID, originAccountAfterTransfer.Balance, testify.Anything).Panic("panic in the process")
		accountRepository.On("UpdateBalance", ctx, destinationAccount.ID, destinationAccountAfterTransfer.Balance, testify.Anything).Return(destinationAccountAfterTransfer, errors.New("error on update destination account balance"))

		transactionHandler := mock.NewTransactionHandlerMock()

		repository := mock.NewRepositoryMock()
		repository.On("BeginTx", ctx).Return(transactionHandler, nil)
		repository.On("CommitTx", transactionHandler).Return(nil)
		repository.On("RollbackTx", transactionHandler).Return(nil)

		transferID := "237d3e7e-2f46-44e7-bf2b-f79721459241"
		createdAt := time.Date(2023, 8, 7, 10, 00, 00, 00, time.UTC)
		transferRepository := mock.NewTransferRepositoryMock()
		transfer, err := entity.NewTransfer(transferID, originAccount, destinationAccount, amount, &createdAt)
		assert.Nil(t, err)

		transferAfterTransaction := *transfer
		transferAfterTransaction.OriginAccount = &originAccountAfterTransfer
		transferAfterTransaction.DestinationAccount = &destinationAccountAfterTransfer

		transferRepository.On("Create", ctx, &transferAfterTransaction, testify.Anything).Return(transferAfterTransaction, nil)

		makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, repository)
		input := usecase.NewMakeTransferUseCaseInput(transferID, originAccount.ID, destinationAccount.ID, amount, &createdAt)

		assert.Panics(t, func() {
			_, _ = makeTransferUseCase.Execute(ctx, input)
		}, "panic in the procaess")

		accountRepository.AssertCalled(t, "FindByID", ctx, originAccount.ID)
		accountRepository.AssertCalled(t, "FindByID", ctx, destinationAccount.ID)
		transferRepository.AssertCalled(t, "Create", ctx, testify.Anything, testify.Anything)
		accountRepository.AssertCalled(t, "UpdateBalance", ctx, originAccount.ID, testify.Anything)
		accountRepository.AssertNotCalled(t, "UpdateBalance", ctx, destinationAccount.ID, testify.Anything)
		repository.AssertCalled(t, "BeginTx", ctx)
		repository.AssertNotCalled(t, "CommitTx", testify.Anything)
		repository.AssertCalled(t, "RollbackTx", testify.Anything)
	})

}
