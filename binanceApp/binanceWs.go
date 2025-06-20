package binanceApp

import (
	"fmt"
	"log"
	"orderFlow/libs/shared"
	"os"
	"os/signal"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
)

func InitBinanceWs() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	for {
		select {
		case <-quit:
			log.Println("Received exit signal, exiting...")
			return
		default:
			websocketStreamClient := binance_connector.NewWebsocketStreamClient(true)
			wsAggTradeHandler := func(event *binance_connector.WsAggTradeEvent) {
				processAggTrade(event)
			}
			errHandler := func(err error) {
				fmt.Println("WebSocket error:", err)
			}
			doneCh, stopCh, err := websocketStreamClient.WsCombinedAggTradeServe(shared.SYMBOLS, wsAggTradeHandler, errHandler)

			if err != nil {
				log.Printf("Binance WebSocket connection error: %v, retrying in 2 seconds...", err)
				select {
				case <-quit:
					log.Println("Received exit signal, exiting...")
					return
				case <-time.After(2 * time.Second):
					continue
				}
			}

			select {
			case <-doneCh:
				log.Println("WebSocket connection closed")
				return
			case <-quit:
				log.Println("Received exit signal, closing WebSocket connection...")
				close(stopCh)
				<-doneCh
				log.Println("WebSocket connection closed gracefully")
				return
			}
		}
	}
}

func processAggTrade(event *binance_connector.WsAggTradeEvent) {
	for _, interval := range shared.INTERVALS {
		aggr := shared.GetAggregator(event.Symbol, interval)
		aggr.ProcessNewAggTrade(event.Symbol, event.IsBuyerMaker, event.Quantity, event.Price, 0)
	}
}
