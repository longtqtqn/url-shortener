package model

import (
	"time"
	"url-shortener/internal/domain"

	"github.com/jinzhu/copier"
	"github.com/uptrace/bun"
)

type UserBunModel struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64      `bun:"id,pk,autoincrement"`
	Email         string     `bun:"email,unique,notnull"`
	APIKEY        string     `bun:"apikey,unique,notnull"`
	DeletedAt     *time.Time `bun:"deleted_at,nullzero,soft_delete"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time `bun:"updated_at,nullzero"`
	Plan          string     `bun:"plan,notnull,default:'free'"`
}

func (m *UserBunModel) ToDomain() *domain.User {
	if m == nil {
		return nil
	}
	var d domain.User
	copier.Copy(&d, m)
	return &d
}

func ToUserBunModel(d *domain.User) *UserBunModel {
	if d == nil {
		return nil
	}
	var m UserBunModel
	copier.Copy(&m, d)
	return &m
}
