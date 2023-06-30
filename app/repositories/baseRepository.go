package repositories

import (
	"context"
	"errors"
	"reflect"

	"gorm.io/gorm"
)

// ‌baseRepository represents the repository for managing users.
type Repository[T any] interface {
	Create(model *T) error
	GetByID(id uint) (*T, error)
	Update(model *T) error
	Delete(model *T) error
}

type BaseRepository[T any] struct {
	db  *gorm.DB
	ctx context.Context
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return BaseRepository[T]{
		db: db,
	}
}

// Create creates a new user.
func (br *BaseRepository[T]) Create(base *T) error {
	if bm, ok := reflect.ValueOf(base).Elem().Interface().(interface {
		BeforeCreate(tx *gorm.DB) error
	}); ok {
		if err := bm.BeforeCreate(br.db); err != nil {
			return err
		}
	}
	return br.db.Create(base).Error
}

// Update updates an existing user.
func (br *BaseRepository[T]) Update(base *T) error {
	if bm, ok := reflect.ValueOf(base).Elem().Interface().(interface {
		BeforeUpdate(tx *gorm.DB) error
	}); ok {
		if err := bm.BeforeUpdate(br.db); err != nil {
			return err
		}
	}
	return br.db.Model(base).Updates(base).Error
}

// GetByID retrieves a user by its ID.
func (br *BaseRepository[T]) GetByID(id uint) (*T, error) {
	var base T
	err := br.db.First(&base, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &base, nil
}

// Delete deletes an existing user.
func (br *BaseRepository[T]) Delete(base *T) error {
	return br.db.Delete(base).Error
}

// GetFirst retrieves the first record that matches the specified conditions.
// It returns gorm.ErrRecordNotFound if no matching record is found.
func (repo *BaseRepository[T]) GetFirst(model T, conditions string, args ...interface{}) error {
	return repo.db.Where(conditions, args...).First(model).Error
}
