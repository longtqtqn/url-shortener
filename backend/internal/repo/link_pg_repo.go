package repo

import (
	"context"
	"url-shortener/backend/internal/domain"
	"url-shortener/backend/internal/repo/model"

	"github.com/uptrace/bun"
)

type LinkPGRepository struct {
	db *bun.DB
}

func NewLinkPGRepository(db *bun.DB) domain.LinkRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &LinkPGRepository{db: db}
}
func (r *LinkPGRepository) CreateLink(ctx context.Context, link *domain.Link) error {
	model := model.ToLinkBunModel(link)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}
	link.ID = model.ID
	return nil
}

func (r *LinkPGRepository) GetLinkByShortCode(ctx context.Context, shortCode string) (*domain.Link, error) {
	var model model.LinkBunModel
	err := r.db.NewSelect().
		Model(&model).
		Where("short_code = ?", shortCode).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *LinkPGRepository) GetListLinksByAPIKeyID(ctx context.Context, apiKeyID int64) ([]*domain.Link, error) {
	var models []model.LinkBunModel
	err := r.db.NewSelect().
		Model(&models).
		Where("apikey_id = ?", apiKeyID).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	links := make([]*domain.Link, 0, len(models))
	for _, m := range models {
		links = append(links, m.ToDomain())
	}
	return links, nil
}

func (r *LinkPGRepository) GetLinkNoAPIKey(ctx context.Context, shortCode string) ([]*domain.Link, error) {
	var models []model.LinkBunModel
	err := r.db.NewSelect().
		Model(&models).
		Where("apikey_id IS NULL").
		Where("short_code = ?", shortCode).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	links := make([]*domain.Link, 0, len(models))
	for _, m := range models {
		links = append(links, m.ToDomain())
	}
	return links, nil
}

func (r *LinkPGRepository) SoftDeleteByShortCode(ctx context.Context, apiKeyID int64, shortCode string) error {
	_, err := r.db.NewUpdate().
		Model((*model.LinkBunModel)(nil)).
		Set("deleted_at = NOW()").
		Where("short_code = ?", shortCode).
		Where("apikey_id = ?", apiKeyID).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

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
