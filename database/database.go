package database

import (
	"fmt"
	"log"
	"os"
	"podcast/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Paginator struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

var db *gorm.DB

func Connection() *gorm.DB {
	return db
}

func Connect() {
	var err error

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	conn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"
	dsn := fmt.Sprintf(conn, host, user, password, dbname, port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database, ", err)
	}

	log.Println("connected successfully to the database")
}

func Migrate() {
	db.Exec(`CREATE TYPE Role AS ENUM (
		'listener',
		'creator'
	);`)

	db.Exec(`CREATE TYPE Visibility AS ENUM (
		'draft',
		'public',
		'protected',
		'archived'
	);`)

	db.Exec(`CREATE TYPE SubscriptionStatus AS ENUM (
		'active',
		'cancelled',
		'trialing'
	);`)

	err := db.AutoMigrate(
		&models.User{},
		&models.Podcast{},
		&models.Episode{},
		&models.Account{},
		&models.Subscription{},
	)

	db.Exec(`CREATE UNIQUE INDEX unique_user_podcast ON subscriptions (user_id, podcast_id);`)

	if err != nil {
		log.Fatal("failed to migrate all database tables", err)
	}

	log.Println("migrated all database tables successfully")
}

func Paginate(p Paginator) *gorm.DB {
	if p.Limit == 0 {
		p.Limit = 10
	}

	if p.Page < 1 {
		p.Page = 1
	}

	return db.Scopes(func(db *gorm.DB) *gorm.DB {
		offset := (p.Page - 1) * p.Limit
		return db.Offset(offset).Limit(p.Limit)
	})
}
