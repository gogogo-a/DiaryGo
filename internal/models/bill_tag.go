package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillTag struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	BillId    uuid.UUID `json:"bill_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Bill      Bill      `json:"bill" gorm:"foreignKey:BillId;references:Id"`
	TagId     uuid.UUID `json:"tag_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tag       Tag       `json:"tag" gorm:"foreignKey:TagId;references:Id"`
}

func (BillTag) TableName() string {
	return "bill_tags"
}

func (b *BillTag) BeforeCreate(tx *gorm.DB) error {
	b.Id = uuid.New()
	return nil
}