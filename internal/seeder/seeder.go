package seeder

import (
	"context"
	"log"
	"url-shortener/internal/domain"

	"github.com/uptrace/bun"
)

func SeedApiKey(ctx context.Context, db *bun.DB) error {
	log.Println("START Seeding API Keys...")

	usersToSeed := []*domain.User{
		{
			APIKEY: "key_for_test@example.com",
			Email:  "test@example.com",
		},
		{
			APIKEY: "key_for_test2@example.com",
			Email:  "test2@example.com",
		},
	}

	for _, user := range usersToSeed {
		count, err := db.NewSelect().
			Model((*domain.User)(nil)).
			Where("email = ?", user.Email).
			Count(ctx)

		if err != nil {
			log.Printf("Error checking for user %s: %v", user.Email, err)
			continue
		}
		if count == 0 {
			_, err := db.NewInsert().Model(user).Exec(ctx)
			if err != nil {
				log.Printf("Error seeding user %s: %v", user.Email, err)
				continue
			}
			log.Printf("Seeded user %s successfully!", user.Email)
		} else {
			log.Printf("User %s already exists, skipping.", user.Email)
		}
	}

	log.Println("Seeding process finished.")
	return nil
}
