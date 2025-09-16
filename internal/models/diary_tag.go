package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiaryTag struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	DiaryId   uuid.UUID `json:"diary_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diary     Diary     `json:"diary" gorm:"foreignKey:DiaryId;references:Id"`
	TagId     uuid.UUID `json:"tag_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tag       Tag       `json:"tag" gorm:"foreignKey:TagId;references:Id"`
}

func (DiaryTag) TableName() string {
	return "diary_tags"
}

func (d *DiaryTag) BeforeCreate(tx *gorm.DB) error {
	d.Id = uuid.New()
	return nil
}