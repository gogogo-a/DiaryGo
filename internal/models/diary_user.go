package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiaryUser struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	DiaryId   uuid.UUID `json:"diary_id" gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diary     Diary     `json:"diary" gorm:"foreignKey:DiaryId;references:Id"`
	UserId    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User      User      `json:"user" gorm:"foreignKey:UserId;references:Id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DiaryUser) TableName() string {
	return "diary_users"
}

func (d *DiaryUser) BeforeCreate(tx *gorm.DB) error {
	if d.Id == uuid.Nil {
		d.Id = uuid.New()
	}
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return nil
}

func (d *DiaryUser) BeforeUpdate(tx *gorm.DB) error {
	d.UpdatedAt = time.Now()
	return nil
}
