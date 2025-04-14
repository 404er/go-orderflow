package binanceApp

import (
	"fmt"
	"log"
	"orderFlow/libs/shared"

	binance_connector "github.com/binance/binance-connector-go"
)

func InitBinanceWs() {
	websocketStreamClient := binance_connector.NewWebsocketStreamClient(true)
	wsAggTradeHandler := func(event *binance_connector.WsAggTradeEvent) {
		processAggTrade(event)
	}
	errHandler := func(err error) {
		fmt.Println("err", err)
	}
	doneCh, _, err := websocketStreamClient.WsCombinedAggTradeServe(shared.SYMBOLS, wsAggTradeHandler, errHandler)

	if err != nil {
		log.Fatal("receive err:", err)
	}
	<-doneCh
}

func processAggTrade(event *binance_connector.WsAggTradeEvent) {
	aggr := shared.GetAggregator(event.Symbol)
	aggr.ProcessNewAggTrade(event.Symbol, event.IsBuyerMaker, event.Quantity, event.Price, 0)
}
