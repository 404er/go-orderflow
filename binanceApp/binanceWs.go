package binanceApp

import (
	"fmt"
	footprint "orderFlow/libs/orderflow"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/gookit/goutil/dump"
)

var Aggregators = map[string]*footprint.OrderFlowAggregator{}

func InitBinanceWs() {
	websocketStreamClient := binance_connector.NewWebsocketStreamClient(false, "wss://stream.binance.com:9443")
	wsAggTradeHandler := func(event *binance_connector.WsAggTradeEvent) {
		processAggTrade(event)
		binance_connector.PrettyPrint(Aggregators)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}

	doneCh, stopCh, err := websocketStreamClient.WsAggTradeServe("BTCUSDT", wsAggTradeHandler, errHandler)

	if err != nil {
		fmt.Println(err)
	}
	go func() {
		time.Sleep(time.Second * 30)
		fmt.Println("stop")
		stopCh <- struct{}{}
	}()
	<-doneCh
}

func processAggTrade(event *binance_connector.WsAggTradeEvent) {
	aggr := getAggregator(event.Symbol, "1m")
	if aggr.ActiveCandle.CloseTimeMs <= event.Time {
		aggr.ProcessCloseCandle()
	}
	aggr.ProcessNewAggTrade(event.Symbol, event.IsBuyerMaker, event.Quantity, event.Price)
	dump.P(aggr.ActiveCandle)
}

func getAggregator(symbol string, interval string) *footprint.OrderFlowAggregator {
	aggregator, ok := Aggregators[symbol]
	if !ok {
		aggregator = &footprint.OrderFlowAggregator{Symbol: symbol, Interval: interval}
		Aggregators[symbol] = aggregator
	}
	return aggregator
}
