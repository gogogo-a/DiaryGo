package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiaryVideo struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	DiaryId   uuid.UUID `json:"diary_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diary     Diary     `json:"diary" gorm:"foreignKey:DiaryId;references:Id"`
	VideoUrl  string    `json:"video_url" gorm:"type:text;not null"`
}

func (DiaryVideo) TableName() string {
	return "diary_videos"
}

func (d *DiaryVideo) BeforeCreate(tx *gorm.DB) error {
	d.Id = uuid.New()
	return nil
}