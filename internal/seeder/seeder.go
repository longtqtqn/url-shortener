package seeder

import (
	"context"
	"log"
	"url-shortener/internal/domain"
	"url-shortener/internal/repo/model"

	"github.com/uptrace/bun"
)

func SeedApiKey(ctx context.Context, db *bun.DB) error {
	log.Println("START Seeding API Keys...")

	usersToSeed := []*domain.User{
		{
			APIKEY: "key1test",
			Email:  "test@example.com",
		},
		{
			APIKEY: "key2test",
			Email:  "test2@example.com",
		},
	}

	for _, user := range usersToSeed {
		userModel := model.ToUserBunModel(user)

		count, err := db.NewSelect().
			Model((*model.UserBunModel)(nil)).
			Where("email = ?", user.Email).
			Count(ctx)

		if err != nil {
			log.Printf("Error checking for user %s: %v", user.Email, err)
			continue
		}

		if count == 0 {
			_, err := db.NewInsert().Model(userModel).Exec(ctx)
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
