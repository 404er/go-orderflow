package database

import (
	"fmt"
	"log"
	orderflow "orderFlow/libs/orderflow"

	"github.com/gookit/goutil/dump"
)

func BatchSaveFootprintCandles(candles []orderflow.FootprintCandle) {
	totalLenth := len(candles)
	if totalLenth == 0 {
		return
	}

	// 逐个保存以确保分片正确
	for _, candle := range candles {
		tableName := fmt.Sprintf("footprint_candles_%s", candle.Symbol)

		db := DB.Table(tableName).Create(&candle)
		if db.Error != nil {
			log.Printf("Error saving candle for %s %s: %v", candle.Symbol, candle.Interval, db.Error)
			continue
		}
		log.Println("saved ", candle.Symbol, candle.Interval)
	}
}

func GetCandles(symbol string, interval string, start int64, end int64) []orderflow.FootprintCandle {
	var candles []orderflow.FootprintCandle
	tableName := fmt.Sprintf("footprint_candles_%s", symbol)
	db := DB.Table(tableName).Where("interval = ? AND \"openTimeMs\" >= ? AND \"closeTimeMs\" <= ?", interval, start, end).Order("\"openTimeMs\" ASC").Find(&candles).Limit(100)
	if db.Error != nil {
		dump.P(db.Error)
		return nil
	}
	return candles
}
