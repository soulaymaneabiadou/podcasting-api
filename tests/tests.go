package tests

import (
	"podcast/database"
	"podcast/models"

	"gorm.io/gorm"
)

func deleteDatabase() {
	database.Connection().
		Unscoped().
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.User{})

	database.Connection().
		Unscoped().
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.Podcast{})
}

func init() {
	database.Connect()
	database.Migrate()

	deleteDatabase()
}

func Teardown() {
	deleteDatabase()
}
