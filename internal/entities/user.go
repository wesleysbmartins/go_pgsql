package entities

import "time"

type User struct {
	Id        int
	Name      string
	Username  string
	Email     string
	Password  string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
