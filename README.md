# Footprint&OrderFlow DataServer
Base on [orderflow](https://github.com/focus1691/orderflow) and fix some bug, currently only supported Binance spot and 1m candle
## How to use

Need Postgres database

```
git clone git@github.com:404er/go-orderflow.git

cd go-orderflow

go mod tidy

go run main.go
```
## Backfill
```
go run main.go -s BTCUSDT -start 2025-04-12 -end 2025-04-14

or

go run main.go -s BTCUSDT -start 2025-04-11

until the day before
```
## ToDoList
- [ ] Add Docker
- [X] Backfill data
- [ ] Support more exchange
- [ ] Support more candle(5m 15m 30m)
