package models

import "time"

type RefreshToken struct {
	Id         int
	UserId     int
	Token      string
	CreateDate time.Time
	ExpDate    time.Time
}
