package main

import (
	"lucassantoss1701/bank/configs"
	"lucassantoss1701/bank/internal/infra/database"
	"lucassantoss1701/bank/internal/infra/database/connection"
	"lucassantoss1701/bank/internal/infra/web"
	"lucassantoss1701/bank/internal/infra/web/webserver"
	"lucassantoss1701/bank/internal/infra/web/webserver/routes"
	"lucassantoss1701/bank/internal/usecase"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	time.Local = time.UTC
	configs.Load()
}

// @title Bank API
// @version 1.0
// @description This API aims to provide resources for common operations that occur within a bank.
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	db := connection.Connect(configs.Get().Database.Type, configs.Get().Database.User, configs.Get().Database.Pass, configs.Get().Database.Host, configs.Get().Database.Port, configs.Get().Database.Name)
	defer db.Close()

	connection.Migrate(db)

	accountRepository := database.NewAccountRepository(db)
	transferRepository := database.NewTransferRepository(db)
	baseRepostiory := database.NewRepository(db)

	webserver := webserver.NewWebServer(configs.Get().Server.Host)

	findAccountUseCase := usecase.NewFindAccountUseCase(accountRepository)
	createAccountUseCase := usecase.NewCreateAccountUseCase(accountRepository)
	findBalanceByAccountUseCase := usecase.NewFindBalanceByAccountUseCase(accountRepository)
	loginUseCase := usecase.NewLoginUseCase(accountRepository)

	webAccountHandler := web.NewWebAccountHandler(createAccountUseCase, findAccountUseCase, findBalanceByAccountUseCase, loginUseCase)

	makeTransferUseCase := usecase.NewMakeTransferUseCase(accountRepository, transferRepository, baseRepostiory)
	findTransfersByAccountUseCase := usecase.NewFindTransfersByAccountUseCase(transferRepository)
	webTransferHandler := web.NewWebTransferHandler(makeTransferUseCase, findTransfersByAccountUseCase)

	routes.HandleAccountRoutes(webserver, webAccountHandler)
	routes.HandleTransferRoutes(webserver, webTransferHandler)

	webserver.Start()
}
