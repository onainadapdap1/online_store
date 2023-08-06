package seeds

import (
	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/driver"
	"github.com/onainadapdap1/online_store/helpers"
	"github.com/onainadapdap1/online_store/models"
)

func seedAdmin(db *gorm.DB) {
	adminRoleUsers := 0
	tx := db.Begin()
	tx.Model(&models.User{}).Where("role = ?", "admin").Count(&adminRoleUsers)
	// if adminRoleUsers == 0 {
		password, _ := helpers.HassPass("password")
		user := models.User{
			FullName: "Admin Golang",
			Email:    "admin7@gmail.com",
			Password: password,
			Role:     "admin",
		}

		tx.Set("gorm:association_autoupdate", false).Debug().Create(&user)

		if tx.Error != nil {
			tx.Rollback()
			print(db.Error)
		}
	// }
	tx.Commit()
}
func Seed() {
	db, _ := driver.ConnectDB()

	seedAdmin(db)
}
