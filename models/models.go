package models

import (
	"time"
)

type User struct {
	Id int64
	Email string
	Password string
	Refresh string
	RefreshIssuedAt time.Time
}
