package main

import (
	"github.com/AhEhIOhYou/etomne/pkg/server"
	"github.com/AhEhIOhYou/etomne/pkg/server/constants"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/logger"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("configs/.env"); err != nil {
		logger.WriteLog(logger.Error, constants.ServerNotFoundEnvFile)
	}
}

func main() {
	server.Start()
}
