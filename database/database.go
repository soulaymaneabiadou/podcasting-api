package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Paginator struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type Sorter struct {
	Column    string
	Ascending bool
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

func Paginate(query *gorm.DB, p Paginator) *gorm.DB {
	if p.Limit == 0 {
		p.Limit = 10
	}

	if p.Page < 1 {
		p.Page = 1
	}

	return query.Scopes(func(db *gorm.DB) *gorm.DB {
		offset := (p.Page - 1) * p.Limit
		return db.Offset(offset).Limit(p.Limit)
	})
}

func Sort(query *gorm.DB, s Sorter) *gorm.DB {
	if s.Column == "" {
		return query
	}

	order := s.Column
	if !s.Ascending {
		order += " DESC"
	}

	return query.Order(order)
}
