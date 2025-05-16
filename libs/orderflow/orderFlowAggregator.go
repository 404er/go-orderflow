package footprint

import (
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderFlowAggregator struct {
	Symbol       string
	Interval     string
	ActiveCandle *FootprintCandle
	mu           sync.RWMutex // 添加互斥锁
}

func (o *OrderFlowAggregator) createNewCandle(startTime int64) {
	var now int64
	if startTime == 0 {
		now = time.Now().UTC().UnixMilli()
	} else {
		now = startTime
	}
	openTimeMs := now
	openTime := time.UnixMilli(openTimeMs).UTC().Format("2006/01/02 15:04:05.000")
	closeTimeMs := getNextMinuteTimestamp(openTimeMs, o.Interval)
	closeTime := time.UnixMilli(closeTimeMs).UTC().Format("2006/01/02 15:04:05.000")

	candle := &FootprintCandle{
		UUID:        uuid.New(),
		OpenTime:    openTime,
		CloseTime:   closeTime,
		OpenTimeMs:  openTimeMs,
		CloseTimeMs: closeTimeMs,
		Interval:    o.Interval,
		Symbol:      o.Symbol,
		Delta:       0,
		Volume:      0,
		AggBid:      0,
		AggAsk:      0,
		Open:        0,
		High:        0,
		Low:         0,
		Close:       0,
		PriceLevels: PriceLevelsMap{},
	}
	o.mu.Lock()
	o.ActiveCandle = candle
	o.mu.Unlock()
}

func (o *OrderFlowAggregator) ProcessCloseCandle() *FootprintCandle {
	o.mu.Lock()
	candle := o.ActiveCandle
	o.ActiveCandle = nil
	o.mu.Unlock()
	return candle
}

func (o *OrderFlowAggregator) ProcessNewAggTrade(symbol string, isBuyerMaker bool, quantity string, price string, startTime int64) {
	o.mu.Lock()
	if o.ActiveCandle == nil {
		o.mu.Unlock()
		o.createNewCandle(startTime)
		o.mu.Lock()
	}

	priceFloat, _ := strconv.ParseFloat(price, 64)
	quantityFloat, _ := strconv.ParseFloat(quantity, 64)

	precisionPrice := decimal.NewFromFloat(priceFloat).Round(2)
	precisionPriceStr := precisionPrice.String()

	if _, exists := o.ActiveCandle.PriceLevels[precisionPriceStr]; !exists {
		o.ActiveCandle.PriceLevels[precisionPriceStr] = PriceLevel{
			VolSumAsk: 0,
			VolSumBid: 0,
		}
	}

	priceLevel := o.ActiveCandle.PriceLevels[precisionPriceStr]

	if o.ActiveCandle.Volume == 0 {
		o.ActiveCandle.Open = priceFloat
		o.ActiveCandle.Low = priceFloat
	}
	o.ActiveCandle.Volume += quantityFloat

	if isBuyerMaker {
		o.ActiveCandle.Delta -= quantityFloat
		o.ActiveCandle.AggAsk += quantityFloat
		priceLevel.VolSumBid += quantityFloat
	} else {
		o.ActiveCandle.Delta += quantityFloat
		o.ActiveCandle.AggBid += quantityFloat
		priceLevel.VolSumAsk += quantityFloat
	}
	o.ActiveCandle.PriceLevels[precisionPriceStr] = priceLevel

	if priceFloat > o.ActiveCandle.High {
		o.ActiveCandle.High = priceFloat
	}

	if priceFloat < o.ActiveCandle.Low {
		o.ActiveCandle.Low = priceFloat
	}
	o.ActiveCandle.Close = priceFloat
	o.mu.Unlock()
}

func getNextMinuteTimestamp(timestamp int64, interval string) int64 {
	// 将时间戳转换为分钟
	minute := timestamp / 60000

	var nextMinute int64
	switch interval {
	case "5m":
		// 计算下一个5分钟整点
		nextMinute = ((minute / 5) + 1) * 5
	case "15m":
		// 计算下一个15分钟整点
		nextMinute = ((minute / 15) + 1) * 15
	case "30m":
		// 计算下一个30分钟整点
		nextMinute = ((minute / 30) + 1) * 30
	default: // 1m
		// 下一个整分钟
		nextMinute = minute + 1
	}

	// 将分钟数转换回时间戳
	return nextMinute * 60000
}
