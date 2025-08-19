package repo

import (
	"context"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/repo/model"
	"url-shortener/internal/usecase"

	"github.com/uptrace/bun"
)

type UserPGRepository struct {
	db *bun.DB
}

func NewUserPGRepository(db *bun.DB) usecase.UserRepository {
	return &UserPGRepository{db: db}
}

// Create implements usecase.UserRepository.
func (r *UserPGRepository) Create(ctx context.Context, user *domain.User) error {
	userModel := model.ToUserBunModel(user)
	_, err := r.db.NewInsert().Model(userModel).Exec(ctx)
	return err
}

// FindByAPIKey implements usecase.UserRepository.
func (r *UserPGRepository) FindByAPIKey(ctx context.Context, apiKey string) (*domain.User, error) {
	var userModel model.UserBunModel
	err := r.db.NewSelect().
		Model(&userModel).
		Join("JOIN apikeys ON apikeys.user_id = user_bun_model.id").
		Where("apikeys.key = ?", apiKey).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return userModel.ToDomain(), nil
}

// FindByID implements usecase.UserRepository.
func (r *UserPGRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	userModel := new(model.UserBunModel)
	err := r.db.NewSelect().Model(userModel).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return userModel.ToDomain(), nil
}

func (r *UserPGRepository) SoftDeleteByID(ctx context.Context, userID int64) error {
	_, err := r.db.NewDelete().
		Model((*model.UserBunModel)(nil)).
		Where("id = ?", userID).
		Exec(ctx)
	return err
}

func (r *UserPGRepository) UpdatePlanAndExpiry(ctx context.Context, userID int64, plan string, expiresAt *time.Time) error {
	q := r.db.NewUpdate().Model((*model.UserBunModel)(nil)).
		Set("plan = ?", plan).
		Where("id = ?", userID)
	if expiresAt == nil {
		q = q.Set("plan_expires_at = NULL")
	} else {
		q = q.Set("plan_expires_at = ?", expiresAt)
	}
	_, err := q.Exec(ctx)
	return err
}

func (r *UserPGRepository) CreateAPIKey(ctx context.Context, userID int64, key string) error {
	api := &model.ApiKeyBunModel{UserID: userID, Key: key}
	_, err := r.db.NewInsert().Model(api).Exec(ctx)
	return err
}
