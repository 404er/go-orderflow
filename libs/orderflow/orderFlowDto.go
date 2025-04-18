package footprint

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type PriceLevel struct {
	VolSumAsk float64 `json:"volSumAsk"`
	VolSumBid float64 `json:"volSumBid"`
}

// PriceLevelsMap 是一个自定义类型，用于处理JSON数据
type PriceLevelsMap map[string]PriceLevel

// Scan 实现了sql.Scanner接口，用于将数据库中的JSON数据转换为Go结构
func (p *PriceLevelsMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	// 解析空值
	if len(bytes) == 0 {
		*p = make(PriceLevelsMap)
		return nil
	}

	// 解析JSON数据
	return json.Unmarshal(bytes, p)
}

type FootprintCandle struct {
	ID          int            `json:"id" gorm:"primaryKey"`
	UUID        uuid.UUID      `json:"uuid" gorm:"column:uuid"`
	OpenTime    string         `json:"openTime" gorm:"column:openTime;type:timestamptz"`
	CloseTime   string         `json:"closeTime" gorm:"column:closeTime;type:timestamptz"`
	OpenTimeMs  int64          `json:"openTimeMs" gorm:"column:openTimeMs;type:bigint"`
	CloseTimeMs int64          `json:"closeTimeMs" gorm:"column:closeTimeMs;type:bigint"`
	Interval    string         `json:"interval" gorm:"column:interval"`
	Symbol      string         `json:"symbol" gorm:"column:symbol"`
	VolumeDelta float64        `json:"volumeDelta" gorm:"column:volumeDelta;type:double precision"`
	Volume      float64        `json:"volume" gorm:"column:volume;type:double precision"`
	AggBid      float64        `json:"aggBid" gorm:"column:aggBid;type:double precision default 0"`
	AggAsk      float64        `json:"aggAsk" gorm:"column:aggAsk;type:double precision default 0"`
	AggTickSize int64          `json:"aggTickSize" gorm:"column:aggTickSize"`
	Open        float64        `json:"open" gorm:"column:open;type:double precision"`
	High        float64        `json:"high" gorm:"column:high;type:double precision"`
	Low         float64        `json:"low" gorm:"column:low;type:double precision"`
	Close       float64        `json:"close" gorm:"column:close;type:double precision"`
	PriceLevels PriceLevelsMap `json:"priceLevels" gorm:"column:priceLevels;type:jsonb"` // 使用PostgreSQL的JSONB类型
}

func (FootprintCandle) TableName() string {
	return "footprint_candles"
}
