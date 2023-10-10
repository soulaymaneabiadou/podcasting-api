package database

import (
	"log"
	"podcast/models"
)

func createTypes() {
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

	// a one to one mapping with the stripe status
	db.Exec(`CREATE TYPE SubscriptionStatus AS ENUM (
		'incomplete',
		'incomplete_expired',
		'trialing',
		'active',
		'past_due',
		'canceled',
		'unpaid'
	);`)

	log.Println("created all types successfully")
}

func dropTypes() {
	db.Exec(`DROP TYPE IF EXISTS Role;`)
	db.Exec(`DROP TYPE IF EXISTS Visibility;`)
	db.Exec(`DROP TYPE IF EXISTS SubscriptionStatus;`)

	log.Println("dropped all types successfully")
}

func createIndices() {
	db.Exec(`CREATE UNIQUE INDEX unique_user_podcast ON subscriptions (user_id, podcast_id);`)

	log.Println("created all indices successfully")
}

func Migrate() {
	createTypes()

	err := db.AutoMigrate(
		&models.User{},
		&models.Podcast{},
		&models.Episode{},
		&models.Subscription{},
	)
	if err != nil {
		log.Fatal("failed to migrate all database tables, ", err)
	}

	createIndices()

	log.Println("migrated all database tables and their indices successfully")
}

func Drop() {
	err := db.Migrator().DropTable(
		&models.User{},
		&models.Podcast{},
		&models.Episode{},
		&models.Subscription{},
	)
	if err != nil {
		log.Fatal("failed to drop all database tables, ", err)
	}

	dropTypes()

	log.Println("dropped all database tables successfully")
}
