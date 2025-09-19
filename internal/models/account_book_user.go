package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountBookUser struct {
	Id            uuid.UUID   `json:"id" gorm:"primaryKey;type:char(36)"`
	AccountBookId uuid.UUID   `json:"account_book_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AccountBook   AccountBook `json:"account_book" gorm:"foreignKey:AccountBookId;references:Id"`
	UserId        uuid.UUID   `json:"user_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User          User        `json:"user" gorm:"foreignKey:UserId;references:Id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (AccountBookUser) TableName() string {
	return "account_book_users"
}

func (a *AccountBookUser) BeforeCreate(tx *gorm.DB) error {
	a.Id = uuid.New()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}
func (a *AccountBookUser) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}
