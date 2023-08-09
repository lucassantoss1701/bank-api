package test

import (
	"bytes"
	"context"
	"encoding/json"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/entity/mock"
	"lucassantoss1701/bank/internal/infra/database"
	"lucassantoss1701/bank/internal/infra/web"
	"lucassantoss1701/bank/internal/infra/web/webserver"
	"lucassantoss1701/bank/internal/infra/web/webserver/routes"
	"lucassantoss1701/bank/internal/usecase"
	"net/http"
	"time"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestCreateAccountHandler(t *testing.T) {
	t.Run("Testing create account with success", func(t *testing.T) {
		UpSQLite3(t)
		defer DownSQLite3(t)

		conn := GetConnection()
		accountRepository := database.NewAccountRepository(conn)

		webAccountHandler := web.NewWebAccountHandler(accountRepository)

		webserver := webserver.NewWebServer(":5000")

		routes.HandleAccountRoutes(webserver, webAccountHandler)

		go webserver.Start()
		time.Sleep(time.Second)
		defer webserver.Stop()

		account := mock.CreateAccount()
		dto := usecase.NewCreateAccountUseCaseInput(account.ID, account.Name, account.CPF, account.Secret, account.Balance, *account.CreatedAt)
		dtoJSON, _ := json.Marshal(dto)

		req, err := http.NewRequest("POST", "http://localhost:5000/accounts", bytes.NewReader(dtoJSON))
		assert.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		defer res.Body.Close()

		var output usecase.CreateAccountUseCaseOutput
		err = json.NewDecoder(res.Body).Decode(&output)
		assert.Nil(t, err)

		query := "SELECT id, name, balance FROM account WHERE id = ?"
		row := conn.QueryRowContext(context.Background(), query, account.ID)

		var accountInDatabase entity.Account
		err = row.Scan(&accountInDatabase.ID, &accountInDatabase.Name, &accountInDatabase.Balance)

		assert.Nil(t, err)
		assert.Equal(t, account.ID, accountInDatabase.ID)
		assert.Equal(t, account.Name, accountInDatabase.Name)
		assert.Equal(t, account.Balance, accountInDatabase.Balance)

	})

}
