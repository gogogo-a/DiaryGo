package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DPermission struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	PermissionName string `json:"permission_name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DPermission) TableName() string {
	return "dpermissions"
}

func (d *DPermission) BeforeCreate(tx *gorm.DB) error {
	d.Id = uuid.New()
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return nil
}

func (d *DPermission) BeforeUpdate(tx *gorm.DB) error {
	d.UpdatedAt = time.Now()
	return nil
}