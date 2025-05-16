package database

import (
	"log"
	orderflow "orderFlow/libs/orderflow"
	"orderFlow/libs/shared"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/sharding"
)

var DB *gorm.DB

func InitDatabaseClient() {
	db, err := gorm.Open(postgres.Open(shared.DB_URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}

	// 配置分库分表
	shardingConfig := sharding.Config{
		ShardingKey:         "Symbol",
		NumberOfShards:      uint(len(shared.SYMBOLS)),
		PrimaryKeyGenerator: sharding.PKSnowflake,
		ShardingSuffixs: func() []string {
			// 生成所有可能的分片后缀
			suffixs := make([]string, len(shared.SYMBOLS))
			for i, symbol := range shared.SYMBOLS {
				suffixs[i] = "_" + symbol
			}
			return suffixs
		},
		ShardingAlgorithm: func(value interface{}) (suffix string, err error) {
			switch v := value.(type) {
			case string:
				return "_" + v, nil
			case int:
				return "_0", nil
			default:
				return "_0", nil
			}
		},
	}

	// 注册分片中间件
	middleware := sharding.Register(shardingConfig, &orderflow.FootprintCandle{})
	db.Use(middleware)

	// 自动迁移表结构
	err = db.AutoMigrate(&orderflow.FootprintCandle{})
	if err != nil {
		panic(err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Client database success")
	DB = db
}
