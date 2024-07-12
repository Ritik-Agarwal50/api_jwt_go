package database

import (
	"fmt"
	"log"

	"github.com/driver/mysql"
	"github.com/joho/gotdotenv"
	"gorm.io/gorm"
)

var GlobalDB *gorm.DB

func InitDatabase() (err error) {
	config, err := gotdotenv.Read()
	if err != nil {
		log.Fatal(err)
	}
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config["DB_USER"], config["DB_PASSWORD"],
		config["DB_HOST"],
		config["DB_NAME"],
	)

	GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return
}
