package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Diary 日记模型
type Diary struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Address   string    `json:"address" gorm:"size:255"`
	Pageview  int       `json:"pageview" gorm:"default:0"`
	Like      int       `json:"like" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

}

// TableName 指定表名
func (Diary) TableName() string {
	return "diaries"
}

// BeforeCreate 创建前的钩子
func (d *Diary) BeforeCreate(tx *gorm.DB) error {
	if d.Id == uuid.Nil {
		d.Id = uuid.New()
	}
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (d *Diary) BeforeUpdate(tx *gorm.DB) error {
	d.UpdatedAt = time.Now()
	return nil
}
