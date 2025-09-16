package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Diary 日记模型
type Diary struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User      User      `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Address   string    `json:"address" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

}

// TableName 指定表名
func (Diary) TableName() string {
	return "diaries"
}

// BeforeCreate 创建前的钩子
func (d *Diary) BeforeCreate(tx *gorm.DB) error {
	d.Id = uuid.New()
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (d *Diary) BeforeUpdate(tx *gorm.DB) error {
	d.UpdatedAt = time.Now()
	return nil
}
