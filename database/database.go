package database

import (
	"fmt"
	"log"
	"os"
	"podcast/models"

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
	DB.Exec(`CREATE TYPE Role AS ENUM (
		'listener',
		'creator'
	);`)

	err := DB.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatal("failed to migrate all database tables", err)
	}

	log.Println("migrated all database tables successfully")
}

type Paginator struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

func Paginate(p Paginator) *gorm.DB {
	if p.Limit == 0 {
		p.Limit = 10
	}

	if p.Page < 1 {
		p.Page = 1
	}

	return DB.Scopes(func(db *gorm.DB) *gorm.DB {
		offset := (p.Page - 1) * p.Limit
		return db.Offset(offset).Limit(p.Limit)
	})
}
