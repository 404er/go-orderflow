package cron

import (
	"log"
	"orderFlow/libs/database"
	"orderFlow/libs/shared"

	"github.com/robfig/cron/v3"
)

func RunTask() {
	go func() {
		c := cron.New()
		quene := database.GetCandleQueneInstance()
		c.AddFunc("*/1 * * * *", func() {
			saveCandleQuene(quene)
		})
		c.Start()
		log.Println("Cron task started")
	}()
}

func saveCandleQuene(quene *database.CandleQuene) {
	aggregators := shared.Aggregators
	for _, aggr := range aggregators {
		candle := aggr.ProcessCloseCandle()
		quene.AddCandle(candle)
		delete(shared.Aggregators, aggr.Symbol)
	}
	quene.Save()
}
