package main

import (
	"github.com/AhEhIOhYou/etomne/backend"
	"github.com/AhEhIOhYou/etomne/backend/constants"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/logger"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("configs/.env"); err != nil {
		logger.WriteLog(logger.Error, constants.ServerNotFoundEnvFile)
	}
}

func main() {
	backend.Start()
}
