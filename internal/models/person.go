package models

import (
	"gorm.io/gorm"
	"time"
)

type Person struct {
	gorm.Model
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Country    string `json:"country"`
}

type PersonSwagger struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Patronymic string     `json:"patronymic"`
	Age        int        `json:"age"`
	Gender     string     `json:"gender"`
	Country    string     `json:"country"`
}
