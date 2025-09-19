package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiaryImage struct {
	Id       uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	DiaryId  uuid.UUID `json:"diary_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diary    Diary     `json:"diary" gorm:"foreignKey:DiaryId;references:Id"`
	ImageUrl string    `json:"image_url" gorm:"type:text;not null"`
}

func (DiaryImage) TableName() string {
	return "diary_images"
}

func (d *DiaryImage) BeforeCreate(tx *gorm.DB) error {
	if d.Id == uuid.Nil {
		d.Id = uuid.New()
	}
	return nil
}
