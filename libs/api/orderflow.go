package api

import (
	"math"
	"orderFlow/libs/database"
	orderflow "orderFlow/libs/orderflow"
	"orderFlow/libs/shared"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	BTCTickSize  = 0.1
	ETHTickSize  = 0.01
	DOGETickSize = 0.0001
	SOLTickSize  = 0.0001
	TickSize     = 10
)

func GetFootprintCandlesService(symbol, interval, stepSize string) []orderflow.FootprintCandle {
	start := time.Now().UnixMilli() - 1000*60*60*24*3
	end := time.Now().UnixMilli()
	if !strings.HasSuffix(symbol, "USDT") {
		return nil
	}
	candles := database.GetCandles(symbol, interval, start, end)
	newCandle := shared.GetActiveCandles(symbol, interval)
	if newCandle != nil {
		candles = append(candles, *newCandle)
	}
	parseCandles(&candles, stepSize)
	return candles
}

func parseCandles(candles *[]orderflow.FootprintCandle, stepSize string) {
	stepSizeInt, _ := strconv.Atoi(stepSize)
	for i := range *candles {
		(*candles)[i].AggBid, _ = decimal.NewFromFloatWithExponent((*candles)[i].AggBid, -2).Float64()
		(*candles)[i].AggAsk, _ = decimal.NewFromFloatWithExponent((*candles)[i].AggAsk, -2).Float64()
		(*candles)[i].Delta, _ = decimal.NewFromFloatWithExponent((*candles)[i].Delta, -2).Float64()
		(*candles)[i].Volume, _ = decimal.NewFromFloatWithExponent((*candles)[i].Volume, -2).Float64()
		switch (*candles)[i].Symbol {
		case "BTCUSDT":
			parseStepSize(&(*candles)[i].PriceLevels, BTCTickSize*float64(TickSize)*float64(stepSizeInt))
		case "ETHUSDT":
			parseStepSize(&(*candles)[i].PriceLevels, ETHTickSize*float64(TickSize)*float64(stepSizeInt))
		case "DOGEUSDT":
			parseStepSize(&(*candles)[i].PriceLevels, DOGETickSize*float64(TickSize)*float64(stepSizeInt))
		case "SOLUSDT":
			parseStepSize(&(*candles)[i].PriceLevels, SOLTickSize*float64(TickSize)*float64(stepSizeInt))
		}
	}
}

func parseStepSize(priceLevels *orderflow.PriceLevelsMap, step float64) {
	newPriceLevels := make(orderflow.PriceLevelsMap)
	// 合并价格级别
	for priceStr, priceLevel := range *priceLevels {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			continue
		}

		// 计算新的合并价格级别
		roundedPrice := math.Floor(price/step) * step
		roundedPriceStr := strconv.FormatFloat(roundedPrice, 'f', -1, 64)

		// 如果已存在该价格级别，则累加数据
		if existingLevel, found := newPriceLevels[roundedPriceStr]; found {
			existingLevel.VolSumAsk += priceLevel.VolSumAsk
			existingLevel.VolSumBid += priceLevel.VolSumBid
			newPriceLevels[roundedPriceStr] = existingLevel
		} else {
			newPriceLevels[roundedPriceStr] = priceLevel
		}
	}
	// 处理交易量数据的精度
	for priceStr, priceLevel := range newPriceLevels {
		// 交易量处理逻辑：小于0.01的取整为0，大于0.01的保留1位小数
		var newVolSumAsk, newVolSumBid float64

		if priceLevel.VolSumAsk < 0.01 {
			newVolSumAsk = 0
		} else if priceLevel.VolSumAsk < 1 && priceLevel.VolSumAsk > 0.01 {
			// 保留1位小数 - 使用 decimal 包确保精度
			askDecimal := decimal.NewFromFloat(priceLevel.VolSumAsk)
			askDecimal = askDecimal.Round(1)
			newVolSumAsk, _ = askDecimal.Float64()
		} else {
			// 保留整数部分
			askDecimal := decimal.NewFromFloat(priceLevel.VolSumAsk)
			askDecimal = askDecimal.Round(0)
			newVolSumAsk, _ = askDecimal.Float64()
		}

		if priceLevel.VolSumBid < 0.01 {
			newVolSumBid = 0
		} else if priceLevel.VolSumBid < 1 && priceLevel.VolSumBid > 0.01 {
			bidDecimal := decimal.NewFromFloat(priceLevel.VolSumBid)
			bidDecimal = bidDecimal.Round(1)
			newVolSumBid, _ = bidDecimal.Float64()
		} else {
			// 保留整数部分
			bidDecimal := decimal.NewFromFloat(priceLevel.VolSumBid)
			bidDecimal = bidDecimal.Round(0)
			newVolSumBid, _ = bidDecimal.Float64()
		}

		newPriceLevels[priceStr] = orderflow.PriceLevel{
			VolSumAsk: newVolSumAsk,
			VolSumBid: newVolSumBid,
		}
	}
	*priceLevels = newPriceLevels
}
