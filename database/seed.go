package database

import (
	"log"
	"time"

	"podcast/models"
)

func Seed() {
	log.Println("Seeding the database with initial data...")

	db := Connection()

	user := models.User{
		Name:       "John Doe",
		Email:      "jdoe@gmail.com",
		Password:   "12345678",
		Role:       "creator",
		Verified:   true,
		VerifiedAt: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		log.Println(err)
	}

	log.Println("seeded the creator")

	podcast := models.Podcast{
		Name:        "Podcastly Gopher",
		Slug:        "podcastly-gopher",
		Headline:    "Vel magnam officiis cupiditate sint sed quia reiciendis.",
		Description: "Ducimus adipisci iusto officia laudantium officiis. In aut id quasi. Quia quis sapiente voluptatem.",
		Picture:     "https://res.cloudinary.com/soulaymaneabiadou/image/upload/v1696504466/Podcastly/Covers/small.jpg",
		SocialLinks: models.SocialLinks{
			Instagram: "https://instagram.com/soulaymaneabiadou",
			Twitter:   "https://twitter.com/soulaymanedev",
		},
		Hosts: []string{
			"John Doe",
			"Adam Smith",
		},
		Tags: []string{
			"dolores",
			"omnis",
			"voluptatem",
		},
		CreatorId: user.ID,
	}

	if err := db.Create(&podcast).Error; err != nil {
		log.Println(err)
	}

	log.Println("seeded the podcast")
}
