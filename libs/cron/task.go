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
			saveCandleQuene(quene, "1m")
		})
		c.AddFunc("*/5 * * * *", func() {
			saveCandleQuene(quene, "5m")
		})
		c.AddFunc("*/15 * * * *", func() {
			saveCandleQuene(quene, "15m")
		})
		c.AddFunc("*/30 * * * *", func() {
			saveCandleQuene(quene, "30m")
		})
		c.Start()
		log.Println("Cron task started")
	}()
}

func saveCandleQuene(quene *database.CandleQuene, interval string) {
	aggregators := shared.Aggregators
	for _, aggr := range aggregators {
		if aggr.Interval != interval {
			continue
		}
		candle := aggr.ProcessCloseCandle()
		quene.AddCandle(candle)
		delete(shared.Aggregators, aggr.Symbol+"_"+aggr.Interval)
	}
	quene.Save()
}
