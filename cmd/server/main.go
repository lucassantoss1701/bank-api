package main

import (
	"database/sql"
	"fmt"
	"lucassantoss1701/bank/configs"
	"lucassantoss1701/bank/internal/infra/database"
	"lucassantoss1701/bank/internal/infra/web"
	"lucassantoss1701/bank/internal/infra/web/webserver"
	"lucassantoss1701/bank/internal/infra/web/webserver/routes"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	time.Local = time.UTC

	configs.Load()
}

type Configs struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func main() {
	configs := Configs{
		DBUser:     "root",
		DBPassword: "root",
		DBHost:     "localhost",
		DBPort:     "3307",
		DBName:     "bank",
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Conexão com o banco de dados MySQL estabelecida com sucesso!")

	accountRepository := database.NewAccountRepository(db)
	transferRepository := database.NewTransferRepository(db)
	baseRepostiory := database.NewRepository(db)

	webserver := webserver.NewWebServer(":8000")
	webAccountHandler := web.NewWebAccountHandler(accountRepository)

	webTransferHandler := web.NewWebTransferHandler(accountRepository, transferRepository, baseRepostiory)

	routes.HandleAccountRoutes(webserver, webAccountHandler)
	routes.HandleTransferRoutes(webserver, webTransferHandler)

	webserver.Start()
}
