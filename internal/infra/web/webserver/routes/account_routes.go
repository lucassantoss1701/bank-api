package routes

import (
	"lucassantoss1701/bank/internal/infra/web"
	"lucassantoss1701/bank/internal/infra/web/webserver"
	"net/http"
)

func HandleAccountRoutes(webserver *webserver.WebServer, webAccountHandler *web.WebAccountHandler) {
	webserver.AddHandler("/accounts", http.MethodGet, webAccountHandler.Find, true)
	webserver.AddHandler("/accounts", http.MethodPost, webAccountHandler.Create, false)
	webserver.AddHandler("/accounts/{account_id}/balance", http.MethodGet, webAccountHandler.FindBalanceByAccount, true)
	webserver.AddHandler("/login", http.MethodPost, webAccountHandler.Login, false)

}
