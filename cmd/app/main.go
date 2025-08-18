package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"url-shortener/internal/repo"
	"url-shortener/internal/seeder"
	"url-shortener/internal/transport/http/handler"
	"url-shortener/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/fx"
)

func requireEnv(keys ...string) {
	missing := []string{}
	for _, k := range keys {
		if _, ok := os.LookupEnv(k); !ok || os.Getenv(k) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		log.Fatalf("Missing required environment variables: %v", missing)
	}
}

func loadEnv() {
	_ = godotenv.Load()
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	_ = godotenv.Overload(".env." + env)
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}
	// Always require essential envs for both development and production
	requireEnv("DATABASE_URL", "PORT", "GIN_MODE", "FREE_PLAN_MAX_LINKS")
}

func main() {
	loadEnv()

	fx.New(
		fx.Provide(
			NewBunDB,
			repo.NewLinkPGRepository,
			repo.NewUserPGRepository,
			usecase.NewShortenerService,
			handler.NewLinkHttpHandler,
		),
		fx.Invoke(RunServer),
	).Run()
}

func NewBunDB() *bun.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:passhihihi@localhost:5432/urlshortener?sslmode=disable" // via pgpool
	}

	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqlDB := sql.OpenDB(connector)
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	db := bun.NewDB(sqlDB, pgdialect.New())

	return db
}

func RunServer(lc fx.Lifecycle, h *handler.LinkHttpHandler, userRepo usecase.UserRepository, db *bun.DB) {
	r := gin.Default()

	r.HEAD("/healthz", func(c *gin.Context) {
		if err := db.RunInTx(c, nil, func(ctx context.Context, tx bun.Tx) error { return nil }); err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.Status(http.StatusOK)
	})

	h.RegisterRoutes(r, userRepo)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Tables and indexes are now managed by migrations.
			// Please run migrations before starting the app.

			//Seed
			if err := seeder.SeedApiKey(ctx, db); err != nil {
				log.Fatalf("Failed to seed user: %v", err)
			}

			//Start server
			go func() {
				addr := os.Getenv("PORT")
				if addr == "" {
					addr = ":8080"
				} else if !strings.HasPrefix(addr, ":") {
					addr = ":" + addr
				}
				log.Println("Server starting on", addr)
				if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Run fail: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
