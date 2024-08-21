package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	Id          string  `json:"id"`
	Username    string  `json:"username"`
	PhoneNumber string  `json:"phone_number"`
	Password    string  `json:"password"`
	Verified    bool    `json:"verified"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
}

func (c UserModel) TableName() string {
	return "users"
}

func (l *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	l.Id = uuid.NewString()
	l.CreatedAt = time.Now().UTC().Format("2006-01-02 15:04:05")
	return
}

func (l *UserModel) BeforeUpdate(tx *gorm.DB) (err error) {
	tNow := time.Now().UTC().Format("2006-01-02 15:04:05")
	l.UpdatedAt = &tNow
	return
}
