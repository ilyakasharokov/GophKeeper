package models

import "time"

type User struct {
	ID       string
	Login    string
	Password string
}

type Note struct {
	ID        string
	Type      int
	Title     string
	Body      []byte
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
