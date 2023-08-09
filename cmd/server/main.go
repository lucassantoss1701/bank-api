package main

import (
	"database/sql"
	"fmt"
	"lucassantoss1701/bank/internal/infra/database"
	"lucassantoss1701/bank/internal/infra/web"
	"lucassantoss1701/bank/internal/infra/web/webserver"
	"lucassantoss1701/bank/internal/infra/web/webserver/routes"

	_ "github.com/go-sql-driver/mysql"
)

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

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
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

	webserver := webserver.NewWebServer(":8000")
	webAccountHandler := web.NewWebAccountHandler(accountRepository)

	routes.HandleAccountRoutes(webserver, webAccountHandler)

	webserver.Start()
}
