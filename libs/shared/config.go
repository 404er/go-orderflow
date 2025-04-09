package shared

import (
	"strings"

	"github.com/gookit/ini/v2"
)

var SYMBOLS []string
var DB_URL string

func InitConfig() {
	err := ini.LoadExists("config.ini")
	if err != nil {
		panic(err)
	}
	SYMBOLS = strings.Split(ini.String("SYMBOLS"), ",")
	DB_URL = ini.String("DB_URL")
}
