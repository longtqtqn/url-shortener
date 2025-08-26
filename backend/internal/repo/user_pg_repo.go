package repo

import (
	"context"
	"time"
	"url-shortener/backend/internal/domain"
	"url-shortener/backend/internal/repo/model"

	"github.com/uptrace/bun"
)

type UserPGRepository struct {
	db *bun.DB
}

func (u *UserPGRepository) CreateUser(ctx context.Context, user *domain.User) error {
	model := model.ToUserBunModel(user)
	_, err := u.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}
	user.ID = model.ID
	user.DeletedAt = nil
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = nil
	return nil
}

// CreateAPIKey implements domain.UserRepository.
func (u *UserPGRepository) CreateAPIKey(ctx context.Context, userID int64, key string) error {
	model := model.ToApiKeyBunModel(&domain.ApiKey{
		UserID: userID,
		Key:    key,
	})
	_, err := u.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}
	model.ID = -1
	model.UserID = userID
	model.Key = key
	model.DeletedAt = nil
	model.CreatedAt = time.Now()
	return nil
}

func (u *UserPGRepository) GetAPIKeysByUserID(ctx context.Context, userID int64, limit, offset int) ([]*domain.ApiKey, error) {
	var models []model.ApiKeyBunModel
	query := u.db.NewSelect().
		Model(&models).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL")

	if limit > 0 {
		query.Limit(limit)
	}
	if offset > 0 {
		query.Offset(offset)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	apiKeys := make([]*domain.ApiKey, 0, len(models))
	for _, m := range models {
		apiKeys = append(apiKeys, m.ToDomain())
	}
	return apiKeys, nil
}

func (u *UserPGRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model model.UserBunModel
	err := u.db.NewSelect().
		Model(&model).
		Where("email = ?", email).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (u *UserPGRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	var model model.UserBunModel
	err := u.db.NewSelect().
		Model(&model).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

// SoftDeleteByID implements domain.UserRepository.
func (u *UserPGRepository) SoftDeleteByID(ctx context.Context, userID int64) error {
	_, err := u.db.NewUpdate().
		Model((*model.UserBunModel)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", userID).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

// GetAPIKeyIDByAPIKey implements domain.UserRepository.
func (u *UserPGRepository) GetAPIKeyIDByAPIKey(ctx context.Context, apiKey string) (int64, error) {
	var apiKeyID int64
	err := u.db.NewSelect().
		Model((*model.ApiKeyBunModel)(nil)).
		Column("id").
		Where("key = ?", apiKey).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(ctx, &apiKeyID)
	if err != nil {
		return 0, err
	}
	return apiKeyID, nil
}

func NewUserPGRepository(db *bun.DB) domain.UserRepository {
	return &UserPGRepository{db: db}
}
