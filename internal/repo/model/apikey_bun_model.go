package model

import (
	"time"

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
