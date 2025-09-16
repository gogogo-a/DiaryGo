package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountBook struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AccountBook) TableName() string {
	return "account_books"
}
func (a *AccountBook) BeforeCreate(tx *gorm.DB) error {
	a.Id = uuid.New()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}

func (a *AccountBook) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}