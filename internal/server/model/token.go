package model

import (
	"time"
)

type Token struct {
	AccessToken string
	UserID      int64
	CreatedAt   time.Time
	EndDateAt   time.Time
}
