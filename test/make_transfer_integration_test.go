package test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestMakeTransferHandler(t *testing.T) {
	t.Run("Testing make transfer with success", func(t *testing.T) {
		UpSQLite3(t)
		defer DownSQLite3(t)

		conn := GetConnection()

		accountRepository := database.NewAccountRepository(conn)
		transferRepository := database.NewTransferRepository(conn)
		baseRepostiory := database.NewRepository(conn)

		webTransferHandler := web.NewWebTransferHandler(accountRepository, transferRepository, baseRepostiory)

		webserver := webserver.NewWebServer(":5000")

		routes.HandleTransferRoutes(webserver, webTransferHandler)

		go webserver.Start()
		time.Sleep(time.Second)
		defer webserver.Stop()

		originAccount := mock.CreateAccount()

		query := "INSERT INTO account (id, name, cpf, secret, balance, created_at) VALUES (?, ?, ?, ?, ?, ?)"
		_, err := conn.ExecContext(context.Background(), query, originAccount.ID, originAccount.Name, originAccount.CPF, originAccount.Secret, originAccount.Balance, originAccount.CreatedAt)
		assert.Nil(t, err)

		destinationID := "fc84682a-3045-4bdf-b91c-10be19f89452"
		destinationAccount := mock.CreateAccount()
		destinationAccount.ID = destinationID

		_, err = conn.ExecContext(context.Background(), query, destinationAccount.ID, destinationAccount.Name, destinationAccount.CPF, destinationAccount.Secret, destinationAccount.Balance, originAccount.CreatedAt)
		assert.Nil(t, err)

		transferData := usecase.MakeTransferUseCaseInput{
			DestinationAccount: usecase.MakeTransferUseCaseAccountInput{
				ID: destinationID,
			},
			Amount: 100,
		}

		transferJSON, _ := json.Marshal(transferData)
		req, err := http.NewRequest("POST", "http://localhost:5000/transfers", bytes.NewBuffer(transferJSON))
		req.Header.Set("Authorization", testToken)
		assert.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		defer res.Body.Close()

		var output usecase.MakeTransferUseCaseOutput
		err = json.NewDecoder(res.Body).Decode(&output)
		assert.Nil(t, err)

		assert.Equal(t, transferData.Amount, output.Amount)
		assert.Equal(t, originAccount.ID, output.OriginAccount.ID)
		assert.Equal(t, destinationAccount.ID, output.DestinationAccount.ID)

	})
}
