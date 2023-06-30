package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (bm *BaseModel) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	bm.CreatedAt = now
	bm.UpdatedAt = now
	return nil
}

// BeforeUpdate sets the updated_at field for the model.
func (bm *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	bm.UpdatedAt = time.Now().UTC()
	return nil
}
