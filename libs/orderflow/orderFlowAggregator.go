package footprint

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/gookit/goutil/timex"
	"github.com/shopspring/decimal"
)

type OrderFlowAggregator struct {
	Symbol       string
	Interval     string
	ActiveCandle *FootprintCandle
}

func (o *OrderFlowAggregator) createNewCandle() {
	openTimeMs := timex.Now()
	openTime := openTimeMs.Format("Y-m-d H:I:S")
	closeTimeMs := getNextMinuteTimestamp(openTimeMs.Timestamp())
	closeTime := timex.FromUnix(closeTimeMs).Format("Y-m-d H:I:S")

	candle := FootprintCandle{
		UUID:        uuid.New(),
		OpenTime:    openTime,
		CloseTime:   closeTime,
		OpenTimeMs:  openTimeMs.Timestamp(),
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
	o.ActiveCandle = &candle
}

func (o *OrderFlowAggregator) ProcessCloseCandle() {}

func (o *OrderFlowAggregator) ProcessNewAggTrade(symbol string, isBuyerMaker bool, quantity string, price string) {
	if o.ActiveCandle == nil {
		o.createNewCandle()
	}

	priceFloat, _ := strconv.ParseFloat(price, 64)
	quantityFloat, _ := strconv.ParseFloat(quantity, 64)

	precisionPrice := decimal.NewFromFloat(priceFloat).Round(2)

	if _, exists := o.ActiveCandle.PriceLevels[precisionPrice.String()]; !exists {
		o.ActiveCandle.PriceLevels[precisionPrice.String()] = PriceLevel{
			VolSumAsk: 0,
			VolSumBid: 0,
		}
	}

	priceLevel := o.ActiveCandle.PriceLevels[precisionPrice.String()]

	if o.ActiveCandle.Volume == 0 {
		o.ActiveCandle.Open = priceFloat
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
	o.ActiveCandle.PriceLevels[precisionPrice.String()] = priceLevel

	if priceFloat > o.ActiveCandle.High {
		o.ActiveCandle.High = priceFloat
	}

	if priceFloat < o.ActiveCandle.Low {
		o.ActiveCandle.Low = priceFloat
	}
	// todo close candle if time is up
}

func getNextMinuteTimestamp(timestamp int64) int64 {
	minute := timestamp / 60
	nextMinute := (minute + 1) * 60
	return nextMinute
}
