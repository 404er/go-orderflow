package main

import (
	"flag"
	"orderFlow/binanceApp"
	"orderFlow/libs/cron"
	"orderFlow/libs/database"
	"orderFlow/libs/shared"
	"os"
)

var (
	s     string
	start string
	end   string
)

func main() {
	parseArgs()
	shared.InitConfig()
	database.InitDatabaseClient()
	if len(os.Args) == 1 {
		cron.RunTask()
		binanceApp.InitBinanceWs()
	} else {
		binanceApp.Backfill(s, start, end)
	}
}

func parseArgs() {
	sPtr := flag.String("s", "", "symbol")
	startPtr := flag.String("start", "", "start time like 2024-01-01")
	endPtr := flag.String("end", "", "end time like 2024-01-02")
	flag.Parse()
	s = *sPtr
	start = *startPtr
	end = *endPtr
}
