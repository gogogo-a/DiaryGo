package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bill struct {
	Id            uuid.UUID   `json:"id" gorm:"primaryKey;type:char(36)"`
	AccountBookId uuid.UUID   `json:"account_book_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AccountBook   AccountBook `json:"account_book" gorm:"foreignKey:AccountBookId;references:Id"`
	UserId        uuid.UUID   `json:"user_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User          User        `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Amount        float64     `json:"amount" gorm:"type:decimal(10,2);not null"`
	Type          string      `json:"type" gorm:"type:varchar(255);not null"` //收入，支出
	Remark        string      `json:"remark" gorm:"type:varchar(255);not null"`
	ImageUrl      string      `json:"image_url" gorm:"type:varchar(255);not null"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (Bill) TableName() string {
	return "bills"
}

func (b *Bill) BeforeCreate(tx *gorm.DB) error {
	if b.Id == uuid.Nil {
		b.Id = uuid.New()
	}
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	return nil
}
func (b *Bill) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
