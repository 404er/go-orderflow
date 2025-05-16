package api

import (
	"log"
	orderflow "orderFlow/libs/orderflow"
	"orderFlow/libs/shared"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitApi() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		cdc := cors.DefaultConfig()
		cdc.AllowOrigins = []string{"http://localhost:3000", "https://rushhub.info"}
		router.Use(cors.New(cdc))
		apiGroup := router.Group("/api")
		apiGroup.GET("/footprint/activeCandles", GetActiveCandles)
		apiGroup.GET("/footprint/candles", GetFootprintCandles)
		log.Println("API server started on port ", shared.API_PORT)
		router.Run(":" + shared.API_PORT)
	}()
}

func GetFootprintCandles(c *gin.Context) {
	symbol := c.Query("symbol")
	interval := c.Query("interval")
	stepSize := c.Query("stepSize")
	candles := GetFootprintCandlesService(symbol, interval, stepSize)
	c.JSON(200, candles)
}

func GetActiveCandles(c *gin.Context) {
	symbol := c.Query("symbol")
	interval := c.Query("interval")
	stepSize := c.Query("stepSize")
	aggr := shared.GetActiveCandles(symbol, interval)
	if aggr != nil {
		var candles []orderflow.FootprintCandle
		candles = append(candles, *aggr)
		parseCandles(&candles, stepSize)
	}
	if aggr == nil {
		c.JSON(404, gin.H{
			"error": "No data found for the specified symbol and interval",
		})
		return
	}
	c.JSON(200, aggr)
}
