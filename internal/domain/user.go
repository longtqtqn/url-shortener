package domain

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64  `bun:"id,pk,autoincrement"`
	Email         string `bun:"email,unique,notnull"`
	APIKEY        string `bun:"apikey,unique,notnull"`
}
