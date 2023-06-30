package db

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _db *gorm.DB

func Initialize() {
	hostname := viper.GetString("mariadb.address")
	user := viper.GetString("mariadb.username")
	pass := viper.GetString("mariadb.password")
	port := viper.GetString("mariadb.port")
	database := viper.GetString("mariadb.database")

	fmt.Printf("DB Address: %v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local\n", user, pass, hostname, port, database)
	var connectionStr = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", user, pass, hostname, port, database)
	_db, _ = gorm.Open(mysql.Open(connectionStr), &gorm.Config{})
}

func GetDB() *gorm.DB {
	if _db == nil {
		Initialize()
	}

	return _db
}

func Ping() error {
	return nil
}