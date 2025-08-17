package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/repo/model"
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
	linkModel := model.ToLinkBunModel(link)
	_, err := r.db.NewInsert().Model(linkModel).ExcludeColumn("id").Exec(ctx)
	return err
}

// FindByShortCode implements usecase.LinkRepository.
func (r *LinkPGRepository) FindByShortCode(ctx context.Context, shortCode string) (*domain.Link, error) {
	linkModel := new(model.LinkBunModel)
	err := r.db.NewSelect().Model(linkModel).Where("short_code = ?", shortCode).Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return linkModel.ToDomain(), nil
}

// ListByUser implements usecase.LinkRepository.
func (r *LinkPGRepository) ListByUser(ctx context.Context, userID int64) ([]*domain.Link, error) {
	linkModels := []*model.LinkBunModel{}

	err := r.db.NewSelect().
		Model(&linkModels).
		Where("user_id = ?", userID).
		Where("deleted_at is NULL").
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	var links []*domain.Link
	for _, lm := range linkModels {
		links = append(links, lm.ToDomain())
	}

	return links, nil
}

// SoftDeleteByShortCode implements usecase.LinkRepository.
func (r *LinkPGRepository) SoftDeleteByShortCode(ctx context.Context, userID int64, shortCode string, deletedAt time.Time) error {
	_, err := r.db.NewDelete().
		Model((*model.LinkBunModel)(nil)).
		Where("user_id = ?", userID).
		Where("short_code = ?", shortCode).
		Exec(ctx)
	return err
}

// TrackClick implements usecase.LinkRepository.
func (r *LinkPGRepository) TrackClick(ctx context.Context, shortCode string) error {
	_, err := r.db.NewUpdate().
		Model((*model.LinkBunModel)(nil)).
		Set("click_count = click_count + 1").
		Set("last_clicked_at = NOW()").
		Where("short_code = ?", shortCode).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *LinkPGRepository) FindLinkCountByUserIDAndLongURL(ctx context.Context, userID int64, longURL string) (int, error) {
	count, err := r.db.NewSelect().
		Model((*model.LinkBunModel)(nil)).
		Where("user_id = ?", userID).
		Where("long_url = ?", longURL).
		Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}
