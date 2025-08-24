package model

import (
	"time"
	"url-shortener/backend/internal/domain"

	"github.com/jinzhu/copier"
	"github.com/uptrace/bun"
)

type ApiKeyBunModel struct {
	bun.BaseModel `bun:"table:apikeys"`
	ID            int64      `bun:"id,pk,autoincrement"`
	UserID        int64      `bun:"user_id,notnull"`
	Key           string     `bun:"key,notnull,unique"`
	DeletedAt     *time.Time `bun:"deleted_at,nullzero,soft_delete"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp"`
}

func (m *ApiKeyBunModel) ToDomain() *domain.ApiKey {
	if m == nil {
		return nil
	}
	var d domain.ApiKey
	copier.Copy(&d, m)
	return &d
}

func ToApiKeyBunModel(d *domain.ApiKey) *ApiKeyBunModel {
	if d == nil {
		return nil
	}
	var m ApiKeyBunModel
	copier.Copy(&m, d)
	return &m
}
