package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"
	domain "url-shortener/internal/domain"
	"url-shortener/internal/usecase"

	"github.com/uptrace/bun"
)

type LinkPGRepository struct {
	db *bun.DB
}

func NewLinkPGRepository(db *bun.DB) usecase.LinkRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &LinkPGRepository{db: db}
}

// Create implements usecase.LinkRepository.
func (r *LinkPGRepository) Create(ctx context.Context, link *domain.Link) error {
	_, err := r.db.NewInsert().Model(link).ExcludeColumn("id").Exec(ctx)
	return err
}

// GetByShortCode implements usecase.LinkRepository.
func (r *LinkPGRepository) GetByShortCode(ctx context.Context, shortCode string) (*domain.Link, error) {
	link := new(domain.Link)
	err := r.db.NewSelect().Model(link).Where("short_code = ?", shortCode).Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return link, nil
}

// ListByUser implements usecase.LinkRepository.
func (l *LinkPGRepository) ListByUser(ctx context.Context, userID int64) ([]*domain.Link, error) {
	panic("unimplemented")
}

// SoftDeleteByShortCode implements usecase.LinkRepository.
func (r *LinkPGRepository) SoftDeleteByShortCode(ctx context.Context, userID int64, shortCode string, deletedAt time.Time) error {
	_, err := r.db.NewDelete().
		Model((*domain.Link)(nil)).
		Where("user_id = ?", userID).
		Where("short_code = ?", shortCode).
		Exec(ctx)
	return err
}

// TrackClick implements usecase.LinkRepository.
func (r *LinkPGRepository) TrackClick(ctx context.Context, shortCode string, at time.Time) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Link)(nil)).
		Set("click_count = click_count + 1").
		Set("last_clicked_at = NOW()").
		Where("short_code = ?", shortCode).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

// func (* Link)
