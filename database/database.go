package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	conn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"
	dsn := fmt.Sprintf(conn, host, user, password, db, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database, ", err)
	}

	log.Println("connected successfully to the database")
}

func Migrate() {
	err := DB.AutoMigrate()

	if err != nil {
		log.Fatal("failed to migrate all database tables", err)
	}

	log.Println("migrated all database tables successfully")
}
