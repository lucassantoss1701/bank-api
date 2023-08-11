package test

import (
	"context"
	"encoding/json"
	"fmt"
	"lucassantoss1701/bank/internal/entity/mock"
	"lucassantoss1701/bank/internal/infra/database"
	"lucassantoss1701/bank/internal/infra/web"
	"lucassantoss1701/bank/internal/infra/web/webserver"
	"lucassantoss1701/bank/internal/infra/web/webserver/routes"
	"lucassantoss1701/bank/internal/usecase"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindBalanceByAccountIDHandler(t *testing.T) {
	t.Run("Testing find balance by account with success", func(t *testing.T) {
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

		query := "INSERT INTO account (id, name, cpf, secret, balance, created_at) VALUES (?, ?, ?, ?, ?, ?)"
		_, err := conn.ExecContext(context.Background(), query, account.ID, account.Name, account.CPF, account.Secret, account.Balance, account.CreatedAt)
		assert.Nil(t, err)

		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:5000/accounts/%s/balance", account.ID), nil)
		assert.Nil(t, err)
		req.Header.Set("Authorization", testToken)

		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var output usecase.FindBalanceByAccountUseCaseOutput
		err = json.NewDecoder(res.Body).Decode(&output)
		assert.Nil(t, err)

		assert.Equal(t, account.Balance, output.Balance)
	})

}
