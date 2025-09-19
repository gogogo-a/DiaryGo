package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiaryLike struct {
	Id        uuid.UUID `gorm:"type:char(36);primary_key;" json:"id"`
	DiaryId   uuid.UUID `json:"diary_id" gorm:"type:char(36);index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diary     Diary     `json:"diary" gorm:"foreignKey:DiaryId;references:Id"`
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User      User      `json:"user" gorm:"foreignKey:UserId;references:Id"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;index;"`
}

func (DiaryLike) TableName() string {
	return "diary_likes"
}

func (d *DiaryLike) BeforeCreate(tx *gorm.DB) error {
	if d.Id == uuid.Nil {
		d.Id = uuid.New()
	}
	d.CreatedAt = time.Now()
	return nil
}
