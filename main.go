package main

import (
	"log"

	"github.com/onainadapdap1/online_store/driver"
	"github.com/onainadapdap1/online_store/router"
)

func main() {
	addr := driver.Config.ServiceHost + ":" + driver.Config.ServicePort

	db, err := driver.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	r := router.Router()
	err = r.Run(addr)
	if err != nil {
		log.Fatal("failed to start the server : ", err)
	}
}
