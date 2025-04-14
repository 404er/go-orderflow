package binanceApp

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"orderFlow/libs/database"
	"orderFlow/libs/shared"
	"orderFlow/libs/utils"
	"os"
	"strconv"
	"strings"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/go-resty/resty/v2"
	"github.com/gookit/goutil/fsutil"
)

func Backfill(s, start, end string) {
	if end == "" {
		end = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	startT, err := time.Parse("2006-01-02", start)
	if err != nil {
		log.Fatalln("error start time", err)
	}
	endT, err := time.Parse("2006-01-02", end)
	if err != nil {
		log.Fatalln("error end time", err)
	}

	diffDays := utils.GetDiffDays(startT, endT)
	var urls []string
	for i := 0; i < diffDays; i++ {
		urls = append(urls, spliceReqUrl(s, startT.AddDate(0, 0, i)))
	}
	downloadFile(urls, s)
}

func spliceReqUrl(symbol string, t time.Time) string {
	baseUrl := "https://data.binance.vision/data/spot/daily/aggTrades/"
	date := t.Format("2006-01-02")
	return fmt.Sprintf("%s%s/%s-aggTrades-%s.zip", baseUrl, symbol, symbol, date)
}

func downloadFile(urls []string, symbol string) {
	client := resty.New()
	downloadDir := "./backfillData"
	client.SetOutputDirectory(downloadDir)
	quene := database.GetCandleQueneInstance()

	for _, url := range urls {
		fileName := strings.Split(url, "aggTrades/")[1]
		simpleName := strings.Split(fileName, "/")[1]
		zipFilePath := downloadDir + "/" + simpleName
		csvFilePath := downloadDir + "/" + strings.Replace(simpleName, ".zip", ".csv", 1)
		if fsutil.FileExists(csvFilePath) && fsutil.FileExists(zipFilePath) {
			log.Println("file already exists", simpleName)
			continue
		}
		log.Println("downloading ", fileName)
		_, err := client.R().SetOutput(simpleName).Get(url)
		if err != nil {
			log.Fatalln("error download file", err)
		}
		fsutil.Unzip(zipFilePath, downloadDir)
		processCsvFile(csvFilePath, symbol)
	}
	quene.Save()
}

func processCsvFile(csvFilePath string, symbol string) []*binance_connector.WsAggTradeEvent {
	csvFile, err := os.ReadFile(csvFilePath)
	if err != nil {
		log.Fatalln("error read file", err)
	}
	cn := csv.NewReader(bytes.NewReader(csvFile))
	cra, _ := cn.ReadAll()
	var tradeRecord []*binance_connector.WsAggTradeEvent
	for i := 0; i < len(cra); i++ {
		event := &binance_connector.WsAggTradeEvent{
			Symbol:   symbol,
			Price:    cra[i][1],
			Quantity: cra[i][2],
			TradeTime: func(t string) int64 {
				// 1744416000057157 -> 1744416000057
				t = t[:13]
				s, _ := strconv.ParseInt(t, 10, 64)
				return s
			}(cra[i][5]),
			IsBuyerMaker: func(b string) bool {
				tb, _ := strconv.ParseBool(b)
				return tb
			}(cra[i][6]),
		}
		agg := shared.GetAggregator(event.Symbol)
		if agg.ActiveCandle != nil && event.TradeTime > agg.ActiveCandle.CloseTimeMs {
			candle := agg.ProcessCloseCandle()
			quene := database.GetCandleQueneInstance()
			quene.AddCandle(candle)
			delete(shared.Aggregators, event.Symbol)
		}
		agg.ProcessNewAggTrade(event.Symbol, event.IsBuyerMaker, event.Quantity, event.Price, event.TradeTime)
	}
	return tradeRecord
}
