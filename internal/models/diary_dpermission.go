package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiaryDPermission struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	DiaryId   uuid.UUID `json:"diary_id" gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diary     Diary     `json:"diary" gorm:"foreignKey:DiaryId;references:Id"`
	DPermissionId uuid.UUID `json:"dpermission_id" gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DPermission DPermission `json:"dpermission" gorm:"foreignKey:DPermissionId;references:Id"`
}

func (DiaryDPermission) TableName() string {
	return "diary_dpermissions"
}

func (d *DiaryDPermission) BeforeCreate(tx *gorm.DB) error {
	if d.Id == uuid.Nil {
		d.Id = uuid.New()
	}
	return nil
}

