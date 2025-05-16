package shared

import (
	footprint "orderFlow/libs/orderflow"
)

var Aggregators = map[string]*footprint.OrderFlowAggregator{}

var INTERVALS = []string{"1m", "5m", "15m", "30m"}

func GetAggregator(symbol string, interval string) *footprint.OrderFlowAggregator {
	symbolInterval := symbol + "_" + interval
	aggregator, ok := Aggregators[symbolInterval]
	if !ok {
		aggregator = &footprint.OrderFlowAggregator{Symbol: symbol, Interval: interval}
		Aggregators[symbolInterval] = aggregator
	}
	return aggregator
}

func GetActiveCandles(symbol string, interval string) *footprint.FootprintCandle {
	symbolInterval := symbol + "_" + interval
	aggregator, ok := Aggregators[symbolInterval]
	if !ok {
		return nil
	}
	return aggregator.ActiveCandle
}
