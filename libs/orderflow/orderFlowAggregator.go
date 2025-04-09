package footprint

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderFlowAggregator struct {
	Symbol       string
	Interval     string
	ActiveCandle *FootprintCandle
}

func (o *OrderFlowAggregator) createNewCandle() {
	now := time.Now().UnixMilli()

	openTimeMs := now
	openTime := time.UnixMilli(openTimeMs).Format("2006/01/02 15:04:05.000")
	closeTimeMs := getNextMinuteTimestamp(openTimeMs)
	closeTime := time.UnixMilli(closeTimeMs).Format("2006/01/02 15:04:05.000")

	candle := &FootprintCandle{
		UUID:        uuid.New(),
		OpenTime:    openTime,
		CloseTime:   closeTime,
		OpenTimeMs:  openTimeMs,
		CloseTimeMs: closeTimeMs,
		Interval:    o.Interval,
		Symbol:      o.Symbol,
		VolumeDelta: 0,
		Volume:      0,
		AggBid:      0,
		AggAsk:      0,
		AggTickSize: 0,
		Open:        0,
		High:        0,
		Low:         0,
		Close:       0,
		PriceLevels: PriceLevelsMap{},
	}
	o.ActiveCandle = candle
}

func (o *OrderFlowAggregator) ProcessCloseCandle() *FootprintCandle {
	candle := o.ActiveCandle
	o.ActiveCandle = nil
	return candle
}

func (o *OrderFlowAggregator) ProcessNewAggTrade(symbol string, isBuyerMaker bool, quantity string, price string) {
	if o.ActiveCandle == nil {
		o.createNewCandle()
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
		o.ActiveCandle.VolumeDelta -= quantityFloat
		o.ActiveCandle.AggAsk += quantityFloat
		priceLevel.VolSumBid += quantityFloat
	} else {
		o.ActiveCandle.VolumeDelta += quantityFloat
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

}

func getNextMinuteTimestamp(timestamp int64) int64 {
	minute := timestamp / 60000
	nextMinute := (minute + 1) * 60000
	return nextMinute
}
