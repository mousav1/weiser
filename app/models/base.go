package models

// import "database/sql"
import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
}

// BeforeCreate is a GORM hook that sets the CreatedAt and UpdatedAt timestamps.
func (bm *BaseModel) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	bm.CreatedAt = now
	bm.UpdatedAt = now
	return nil
}

// BeforeUpdate is a GORM hook that sets the UpdatedAt timestamp.
func (bm *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	bm.UpdatedAt = time.Now().UTC()
	return nil
}
