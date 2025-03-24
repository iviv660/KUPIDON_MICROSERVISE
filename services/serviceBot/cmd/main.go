package main

import (
	clientsMatch "serviceBot/internal/clients/match_client"
	clientsUser "serviceBot/internal/clients/user_client"
	"serviceBot/internal/config"
	"serviceBot/internal/usecase"
)

func main() {
	cfg := config.NewConfig()
	serviceUser := clientsUser.NewHTTPUserServiseClient(cfg.USER_SERVICE)
	serviceMatch := clientsMatch.NewHTTPMatchServiseClient(cfg.MATCH_SERVICE)
	uc := usecase.NewUseCase(serviceUser, serviceMatch)
	uc.StartBot(cfg.TELEGRAM_BOT_TOKEN)
}
