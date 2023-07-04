package repositories

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Repository represents the repository for managing models.
type Repository[T any] interface {
	Create(model *T) error
	GetByID(id uint) (*T, error)
	Update(model *T) error
	Delete(model *T) error
	GetFirst(model T, conditions string, args ...interface{}) error
}

// BaseRepository represents the base repository for managing models.
type BaseRepository[T any] struct {
	db  *gorm.DB
	ctx context.Context
}

// NewBaseRepository creates a new instance of the base repository.
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return BaseRepository[T]{
		db: db,
	}
}

// Create creates a new record.
func (br *BaseRepository[T]) Create(base *T) error {
	if base == nil {
		return errors.New("nil model provided")
	}
	if bm, ok := reflect.ValueOf(base).Elem().Interface().(interface {
		BeforeCreate(tx *gorm.DB) error
	}); ok {
		if err := bm.BeforeCreate(br.db); err != nil {
			return errors.Wrap(err, "error in BeforeCreate")
		}
	}
	if err := br.db.Create(base).Error; err != nil {
		return errors.Wrap(err, "error in Create")
	}
	return nil
}

// Update updates an existing record.
func (br *BaseRepository[T]) Update(base *T) error {
	if base == nil {
		return errors.New("nil model provided")
	}
	if bm, ok := reflect.ValueOf(base).Elem().Interface().(interface {
		BeforeUpdate(tx *gorm.DB) error
	}); ok {
		if err := bm.BeforeUpdate(br.db); err != nil {
			return errors.Wrap(err, "error in BeforeUpdate")
		}
	}
	if err := br.db.Model(base).Updates(base).Error; err != nil {
		return errors.Wrap(err, "error in Update")
	}
	return nil
}

// GetByID retrieves a record by its ID.
func (br *BaseRepository[T]) GetByID(id uint) (*T, error) {
	var base T
	if err := br.db.First(&base, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "error in GetByID")
	}
	return &base, nil
}

// Delete deletes an existing record.
func (br *BaseRepository[T]) Delete(base *T) error {
	if base == nil {
		return errors.New("nil model provided")
	}
	if err := br.db.Delete(base).Error; err != nil {
		return errors.Wrap(err, "error in Delete")
	}
	return nil
}

// GetFirst retrieves the first record that matches the specified conditions.
// It returns gorm.ErrRecordNotFound if no matching record is found.
func (br *BaseRepository[T]) GetFirst(model *T, conditions string, args ...interface{}) error {
	if model == nil {
		return errors.New("nil model provided")
	}
	if err := br.db.Where(conditions, args...).First(model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "record not found in GetFirst")
		}
		return errors.Wrap(err, "error in GetFirst")
	}
	return nil
}
