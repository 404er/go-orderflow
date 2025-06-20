package database

import (
	"log"
	orderflow "orderFlow/libs/orderflow"
	shared "orderFlow/libs/shared"
	"sync"
)

var (
	instance *CandleQuene
	once     sync.Once
)

type CandleQuene struct {
	Candles []orderflow.FootprintCandle
	mu      sync.Mutex
}

func GetCandleQueneInstance() *CandleQuene {
	once.Do(func() {
		instance = &CandleQuene{
			Candles: []orderflow.FootprintCandle{},
		}
	})
	return instance
}

func (q *CandleQuene) AddCandle(candle *orderflow.FootprintCandle) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Candles = append(q.Candles, *candle)
}

func (q *CandleQuene) Save() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.Candles) == 0 {
		return
	}
	if !shared.IS_TEST {
		BatchSaveFootprintCandles(q.Candles)
	} else {
		log.Println("test mode, not saving candles")
	}
	q.Candles = []orderflow.FootprintCandle{}
}
