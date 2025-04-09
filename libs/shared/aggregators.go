package shared

import (
	footprint "orderFlow/libs/orderflow"
)

var Aggregators = map[string]*footprint.OrderFlowAggregator{}

func GetAggregator(symbol string) *footprint.OrderFlowAggregator {
	aggregator, ok := Aggregators[symbol]
	if !ok {
		aggregator = &footprint.OrderFlowAggregator{Symbol: symbol, Interval: "1m"}
		Aggregators[symbol] = aggregator
	}
	return aggregator
}
