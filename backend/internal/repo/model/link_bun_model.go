package model

import (
	"time"
	"url-shortener/backend/internal/domain"

	"github.com/jinzhu/copier"
	"github.com/uptrace/bun"
)

type LinkBunModel struct {
	bun.BaseModel `bun:"table:links"`
	ID            int64      `bun:"id,pk,autoincrement"`
	APIKeyID      *int64     `bun:"apikey_id,nullzero"` // Optional FK
	ShortCode     string     `bun:"short_code,notnull,unique"`
	LongURL       string     `bun:"long_url,notnull"`
	Password      *string    `bun:"password,nullzero"` // Optional password
	ClickCount    int64      `bun:"click_count,notnull,default:0"`
	LastClickedAt *time.Time `bun:"last_clicked_at,nullzero"`
	DeletedAt     *time.Time `bun:"deleted_at,nullzero,soft_delete"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp"`
}

func (m *LinkBunModel) ToDomain() *domain.Link {
	if m == nil {
		return nil
	}
	var d domain.Link
	copier.Copy(&d, m)
	return &d
}

func ToLinkBunModel(d *domain.Link) *LinkBunModel {
	if d == nil {
		return nil
	}
	var m LinkBunModel
	copier.Copy(&m, d)
	return &m
}
