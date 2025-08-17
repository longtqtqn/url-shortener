package model

import (
	"time"
	"url-shortener/internal/domain"

	"github.com/uptrace/bun"
)

type LinkBunModel struct {
	bun.BaseModel `bun:"table:links"`
	ID            int64      `bun:"id,pk,autoincrement"`
	UserID        int64      `bun:"user_id,notnull"`
	ShortCode     string     `bun:"short_code,notnull,unique"`
	LongURL       string     `bun:"long_url,notnull"`
	ClickCount    int64      `bun:"click_count,notnull,default:0"`
	LastClickedAt *time.Time `bun:"last_clicked_at,nullzero"`
	DeletedAt     *time.Time `bun:"deleted_at,nullzero,soft_delete"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp"`
}

func (m *LinkBunModel) ToDomain() *domain.Link {
	if m == nil {
		return nil
	}
	return &domain.Link{
		ID:            m.ID,
		UserID:        m.UserID,
		ShortCode:     m.ShortCode,
		LongURL:       m.LongURL,
		ClickCount:    m.ClickCount,
		LastClickedAt: m.LastClickedAt,
		DeletedAt:     m.DeletedAt,
		CreatedAt:     m.CreatedAt,
	}
}

func ToLinkBunModel(link *domain.Link) *LinkBunModel {
	if link == nil {
		return nil
	}
	return &LinkBunModel{
		ID:            link.ID,
		UserID:        link.UserID,
		ShortCode:     link.ShortCode,
		LongURL:       link.LongURL,
		ClickCount:    link.ClickCount,
		LastClickedAt: link.LastClickedAt,
		DeletedAt:     link.DeletedAt,
		CreatedAt:     link.CreatedAt,
	}
}
