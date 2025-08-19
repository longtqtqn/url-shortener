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

type seedUser struct {
	Email  string `json:"email"`
	APIKey string `json:"apikey"`
	Plan   string `json:"plan"`
	Role   string `json:"role"`
}

func parseSeedUsers() []seedUser {
	// Priority: SEED_USERS_JSON > defaults
	if raw := os.Getenv("SEED_USERS_JSON"); strings.TrimSpace(raw) != "" {
		var users []seedUser
		if err := json.Unmarshal([]byte(raw), &users); err == nil && len(users) > 0 {
			return users
		} else if err != nil {
			log.Printf("Failed to parse SEED_USERS_JSON: %v", err)
		}
	}
	// defaults
	return []seedUser{
		{Email: "test@example.com", APIKey: "key1test", Plan: "free", Role: "user"},
		{Email: "test2@example.com", APIKey: "key2test", Plan: "premium", Role: "admin"},
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

	seedUsers := parseSeedUsers()

	for _, su := range seedUsers {
		// Upsert user row first
		user := &domain.User{Email: su.Email, Plan: su.Plan, Role: su.Role}
		userModel := model.ToUserBunModel(user)

		if mode == "exist-only" {
			_, err := db.NewInsert().
				Model(userModel).
				On("CONFLICT (email) DO NOTHING").
				Exec(ctx)
			if err != nil {
				log.Printf("Error inserting user %s: %v", su.Email, err)
				continue
			}
		} else {
			_, err := db.NewInsert().
				Model(userModel).
				On("CONFLICT (email) DO UPDATE").
				Set("plan = EXCLUDED.plan").
				Set("role = EXCLUDED.role").
				Set("deleted_at = NULL").
				Exec(ctx)
			if err != nil {
				log.Printf("Error upserting user %s: %v", su.Email, err)
				continue
			}
		}

		// Ensure we have the user id
		var persistedUser model.UserBunModel
		if err := db.NewSelect().Model(&persistedUser).Where("email = ?", su.Email).Scan(ctx); err != nil {
			log.Printf("Failed to load user id for %s: %v", su.Email, err)
			continue
		}

		// Upsert API key row mapping to user
		apiKey := &model.ApiKeyBunModel{UserID: persistedUser.ID, Key: su.APIKey}
		if mode == "exist-only" {
			_, err := db.NewInsert().Model(apiKey).On("CONFLICT (key) DO NOTHING").Exec(ctx)
			if err != nil {
				log.Printf("Error inserting apikey for %s: %v", su.Email, err)
				continue
			}
			log.Printf("Inserted apikey for %s (or existed)", su.Email)
			continue
		}
		_, err := db.NewInsert().
			Model(apiKey).
			On("CONFLICT (key) DO UPDATE").
			Set("user_id = EXCLUDED.user_id").
			Set("deleted_at = NULL").
			Exec(ctx)
		if err != nil {
			log.Printf("Error upserting apikey for %s: %v", su.Email, err)
			continue
		}
		log.Printf("Upserted apikey for %s successfully!", su.Email)
	}

	log.Println("Seeding process finished.")
	return nil
}
