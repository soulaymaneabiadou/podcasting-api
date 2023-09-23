package tests

import (
	"podcast/database"
	"podcast/models"

	"gorm.io/gorm"
)

func deleteDatabase() {
	database.DB.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})
}

func init() {
	database.Connect()
	database.Migrate()

	deleteDatabase()
}

func Teardown() {
	deleteDatabase()
}
