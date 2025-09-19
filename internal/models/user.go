package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id        uuid.UUID  `json:"id" gorm:"primaryKey;type:char(36)"`
	PlantId   string     `json:"plant_id" gorm:"not null"`
	PlantForm string     `json:"plant_form" gorm:"not null"`
	UserName  string     `json:"user_name" gorm:"not null"` //默认为用户+uuid
	Password  string     `json:"password"`
	Avatar    string     `json:"avatar"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Gender    string     `json:"gender"`
	Birthday  *time.Time `json:"birthday"`
	Address   string     `json:"address"`
	Remark    string     `json:"remark"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Id == uuid.Nil {
		u.Id = uuid.New()
	}
	u.UserName = "用户" + u.Id.String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Birthday = nil
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
