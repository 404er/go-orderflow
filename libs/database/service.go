package database

import (
	"log"
	orderflow "orderFlow/libs/orderflow"
)

func BatchSaveFootprintCandles(candles []orderflow.FootprintCandle) {
	totalLenth := len(candles)
	if totalLenth == 0 {
		return
	}
	db := DB.Create(&candles)
	if db.Error != nil {
		log.Fatal(db.Error)
	}
	log.Println("saved ", totalLenth, " candles")
}

func GetCandles(symbol string, interval string, start int64, end int64) []orderflow.FootprintCandle {
	var candles []orderflow.FootprintCandle
	db := DB.Where("symbol = ? AND interval = ? AND openTimeMs >= ? AND openTimeMs <= ?", symbol, interval, start, end).Order("openTimeMs ASC").Find(&candles)
	if db.Error != nil {
		log.Fatal("Get Candle Error")
		return nil
	}
	return candles
}
