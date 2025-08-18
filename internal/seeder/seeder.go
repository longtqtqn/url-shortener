package seeder

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"url-shortener/internal/domain"
	"url-shortener/internal/repo/model"

	"github.com/uptrace/bun"
)

func parseSeedUsers() []*domain.User {
	// Priority: SEED_USERS_JSON > defaults
	if raw := os.Getenv("SEED_USERS_JSON"); strings.TrimSpace(raw) != "" {
		var users []*domain.User
		if err := json.Unmarshal([]byte(raw), &users); err == nil && len(users) > 0 {
			return users
		} else if err != nil {
			log.Printf("Failed to parse SEED_USERS_JSON: %v", err)
		}
	}
	// defaults
	return []*domain.User{
		{APIKEY: "key1test", Email: "test@example.com", Plan: "free"},
		{APIKEY: "key2test", Email: "test2@example.com", Plan: "premium"},
	}
}

func SeedApiKey(ctx context.Context, db *bun.DB) error {
	if strings.ToLower(os.Getenv("SEED_ENABLED")) != "true" {
		log.Println("Seeding disabled (SEED_ENABLED != true), skipping.")
		return nil
	}
	mode := strings.ToLower(os.Getenv("SEED_MODE"))
	if mode != "exist-only" {
		mode = "enforce"
	}

	log.Println("START Seeding API Keys... (mode:", mode, ")")

	usersToSeed := parseSeedUsers()

	for _, user := range usersToSeed {
		userModel := model.ToUserBunModel(user)

		if mode == "exist-only" {
			// Insert if not exists; do not modify existing rows
			_, err := db.NewInsert().
				Model(userModel).
				On("CONFLICT (email) DO NOTHING").
				Exec(ctx)
			if err != nil {
				log.Printf("Error inserting user %s: %v", user.Email, err)
				continue
			}
			log.Printf("Inserted (or existed) user %s", user.Email)
			continue
		}

		// enforce: upsert and restore if soft-deleted
		_, err := db.NewInsert().
			Model(userModel).
			On("CONFLICT (email) DO UPDATE").
			Set("apikey = EXCLUDED.apikey").
			Set("plan = EXCLUDED.plan").
			Set("deleted_at = NULL").
			Exec(ctx)
		if err != nil {
			log.Printf("Error upserting user %s: %v", user.Email, err)
			continue
		}
		log.Printf("Upserted user %s successfully!", user.Email)
	}

	log.Println("Seeding process finished.")
	return nil
}
