package driver

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	DB  *gorm.DB
	err error
)


func ConnectDB() (*gorm.DB, error) {

	if DB != nil {
		return DB, nil
	}

	dbConfig := Config.DB
	if dbConfig.Adapter == "postgresonai" {
		DB, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbConfig.UserDB, dbConfig.Password, dbConfig.Host, dbConfig.Name))
		log.Println("Connected to Database Local")
	}

	if err != nil {
		log.Println("[Driver.ConnectDB] error when connect to database")
		log.Fatal(err)
	} else {
		log.Println("SUCCESS CONNECT TO DATABASE")
	}

	return DB, nil
}