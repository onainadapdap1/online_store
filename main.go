package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/driver"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/router"
	"github.com/onainadapdap1/online_store/seeds"
)

func initTable(db *gorm.DB) {
	db.Debug().AutoMigrate(&models.User{}).AddUniqueIndex("idx_users_email", "email")
	db.Debug().AutoMigrate(&models.Category{})
	db.Debug().AutoMigrate(&models.Product{})
	db.Debug().AutoMigrate(&models.PaymentCategory{})
	db.Debug().AutoMigrate(&models.PaymentMethod{})
}

func drop(db *gorm.DB) {
	db.DropTableIfExists(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.PaymentCategory{},
		&models.PaymentMethod{},
	)
}

func create(database *gorm.DB) {
	drop(database)
	initTable(database)
}

func main() {
	addr := driver.Config.ServiceHost + ":" + driver.Config.ServicePort

	db, err := driver.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	args := os.Args
	if len(args) > 1 {
		first := args[1]

		if first == "create" {
			create(db)
		}else if first == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			initTable(db)
		}

		if first != "" {
			os.Exit(0)
		}
	}

	r := router.Router()
	err = r.Run(addr)
	if err != nil {
		log.Fatal("failed to start the server : ", err)
	}
}
