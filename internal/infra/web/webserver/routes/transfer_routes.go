package routes

import (
	"lucassantoss1701/bank/internal/infra/web"
	"lucassantoss1701/bank/internal/infra/web/webserver"
	"net/http"
)

func HandleTransferRoutes(webserver *webserver.WebServer, webTransferHandler *web.WebTransferHandler) {
	webserver.AddHandler("/transfers", http.MethodPost, webTransferHandler.Create, true)
	webserver.AddHandler("/transfers", http.MethodGet, webTransferHandler.FindByAccountID, true)

}
