package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"url-shortener/internal/domain"
	"url-shortener/internal/repo"
	"url-shortener/internal/seeder"
	"url-shortener/internal/transport/http/handler"
	"url-shortener/internal/transport/middleware"
	"url-shortener/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/fx"
)

func main() {
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
	dsn := "postgres://user:passhihihi@localhost:5432/urlshortener?sslmode=disable"

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

	api := r.Group("/api")
	api.Use(middleware.ApiKeyAuth(userRepo))
	h.RegisterRoutes(api)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			//Create tables
			if _, err := db.NewCreateTable().Model((*domain.User)(nil)).IfNotExists().Exec(ctx); err != nil {
				log.Fatalf("Failed to create users table: %v", err)
			}
			if _, err := db.NewCreateTable().Model((*domain.Link)(nil)).IfNotExists().Exec(ctx); err != nil {
				log.Fatalf("Failed to create links table: %v", err)
			}

			//Seed
			if err := seeder.SeedApiKey(ctx, db); err != nil {
				log.Fatalf("Failed to seed user: %v", err)
			}

			//Start server
			go func() {
				log.Println("Server starting on :8080")
				if err := r.Run(":8080"); err != nil && err != http.ErrServerClosed {
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
