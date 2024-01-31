package models

import "time"

type RefreshToken struct {
	Id         int64
	UserId     int64
	Token      string
	CreateDate time.Time
	ExpDate    time.Time
}
