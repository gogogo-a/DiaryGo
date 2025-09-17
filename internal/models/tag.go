package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	TagName   string    `json:"tag_name" gorm:"type:varchar(255);not null"`
	Type      string    `json:"type" gorm:"type:varchar(255);not null"`//旅游，餐饮，购物，交通，住宿，其他
	Category  string    `json:"category" gorm:"type:varchar(255);not null"`//账单，日记
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Tag) TableName() string {
	return "tags"
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	t.Id = uuid.New()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Tag) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
