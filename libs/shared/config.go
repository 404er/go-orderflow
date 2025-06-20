package shared

import (
	"log"
	"strings"

	"github.com/gookit/ini/v2"
)

var SYMBOLS []string
var DB_URL string
var API_PORT string
var IS_TEST bool

func InitConfig() {
	err := ini.LoadExists("config.ini")
	if err != nil {
		panic(err)
	}
	SYMBOLS = strings.Split(ini.String("SYMBOLS"), ",")
	DB_URL = ini.String("DB_URL")
	API_PORT = ini.String("API_PORT")
	IS_TEST = ini.Bool("IS_TEST")
	if IS_TEST {
		log.Println("test mode only output active candles")
	}
}
