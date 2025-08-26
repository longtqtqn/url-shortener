package seeder

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"url-shortener/backend/internal/domain"
	"url-shortener/backend/internal/repo/model"

	"github.com/uptrace/bun"
)

type seedUser struct {
	Email  string `json:"email"`
	APIKey string `json:"apikey"`
}

func parseSeedUsers() []seedUser {
	// Priority: SEED_USERS_JSON > defaults
	if raw := os.Getenv("SEED_USERS_JSON"); strings.TrimSpace(raw) != "" {
		var users []seedUser
		if err := json.Unmarshal([]byte(raw), &users); err == nil && len(users) > 0 {
			log.Printf("Loaded %d users from SEED_USERS_JSON", len(users))
			return users
		} else if err != nil {
			log.Printf("Failed to parse SEED_USERS_JSON: %v", err)
		}
	}

	// defaults
	defaultUsers := []seedUser{
		{Email: "test@example.com", APIKey: "testkey12345678901234567890123456"},
		{Email: "demo@example.com", APIKey: "demokey1234567890123456789012345"},
	}
	log.Printf("Using %d default seed users", len(defaultUsers))
	return defaultUsers
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

	log.Printf("START Seeding API Keys... (mode: %s)", mode)

	seedUsers := parseSeedUsers()
	successCount := 0
	errorCount := 0

	for _, su := range seedUsers {
		log.Printf("Processing user: %s", su.Email)

		// Validate API key length (should be 32 characters for hex)
		if len(su.APIKey) != 32 {
			log.Printf("Warning: API key for %s is not 32 characters long (got %d)", su.Email, len(su.APIKey))
		}

		// Upsert user row first
		user := &domain.User{Email: su.Email}
		userModel := model.ToUserBunModel(user)

		if mode == "exist-only" {
			_, err := db.NewInsert().
				Model(userModel).
				On("CONFLICT (email) DO NOTHING").
				Exec(ctx)
			if err != nil {
				log.Printf("Error inserting user %s: %v", su.Email, err)
				errorCount++
				continue
			}
		} else {
			_, err := db.NewInsert().
				Model(userModel).
				On("CONFLICT (email) DO UPDATE").
				Set("deleted_at = NULL").
				Exec(ctx)
			if err != nil {
				log.Printf("Error upserting user %s: %v", su.Email, err)
				errorCount++
				continue
			}
		}

		// Ensure we have the user id
		var persistedUser model.UserBunModel
		if err := db.NewSelect().Model(&persistedUser).Where("email = ?", su.Email).Scan(ctx); err != nil {
			log.Printf("Failed to load user id for %s: %v", su.Email, err)
			errorCount++
			continue
		}

		// Upsert API key row mapping to user
		apiKey := &model.ApiKeyBunModel{UserID: persistedUser.ID, Key: su.APIKey}
		if mode == "exist-only" {
			_, err := db.NewInsert().Model(apiKey).On("CONFLICT (key) DO NOTHING").Exec(ctx)
			if err != nil {
				log.Printf("Error inserting apikey for %s: %v", su.Email, err)
				errorCount++
				continue
			}
			log.Printf("Inserted apikey for %s (or existed)", su.Email)
		} else {
			_, err := db.NewInsert().
				Model(apiKey).
				On("CONFLICT (key) DO UPDATE").
				Set("user_id = EXCLUDED.user_id").
				Set("deleted_at = NULL").
				Exec(ctx)
			if err != nil {
				log.Printf("Error upserting apikey for %s: %v", su.Email, err)
				errorCount++
				continue
			}
			log.Printf("Upserted apikey for %s successfully!", su.Email)
		}

		successCount++
	}

	log.Printf("Seeding process finished. Success: %d, Errors: %d", successCount, errorCount)

	if errorCount > 0 {
		log.Printf("Warning: %d errors occurred during seeding", errorCount)
	}

	return nil
}
