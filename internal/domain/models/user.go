package models

import "time"

type User struct {
	Id           int64
	Email        string
	Name         string
	PasswordHash []byte
	CreateDate   time.Time
}

type NewUser struct {
	Email    string
	Name     string
	Password string
}

type UserResponse struct {
	UserId       int64
	RefreshToken string
	AccessToken  string
}
