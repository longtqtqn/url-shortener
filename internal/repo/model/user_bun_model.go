package model

import (
	"url-shortener/internal/domain"

	"github.com/uptrace/bun"
)

type UserBunModel struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64  `bun:"id,pk,autoincrement"`
	Email         string `bun:"email,unique,notnull"`
	APIKEY        string `bun:"apikey,unique,notnull"`
}

func (m *UserBunModel) ToDomain() *domain.User {
	if m == nil {
		return nil
	}
	return &domain.User{
		ID:     m.ID,
		Email:  m.Email,
		APIKEY: m.APIKEY,
	}
}

func ToUserBunModel(d *domain.User) *UserBunModel {
	if d == nil {
		return nil
	}
	return &UserBunModel{
		ID:     d.ID,
		Email:  d.Email,
		APIKEY: d.APIKEY,
	}
}
