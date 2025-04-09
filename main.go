package main

import (
	"orderFlow/binanceApp"
	"orderFlow/libs/cron"
	"orderFlow/libs/database"
	"orderFlow/libs/shared"
)

func main() {
	shared.InitConfig()
	database.InitDatabaseClient()
	cron.RunTask()
	binanceApp.InitBinanceWs()
}
