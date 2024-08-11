package storage

import (
	"fmt"
)

// Registry برای نگهداری driverها
var Registry = make(map[string]Storage)

// DefaultDriverName نام درایور پیش‌فرض
var DefaultDriverName string

// DefaultDriver متغیر برای ذخیره درایور پیش‌فرض
var DefaultDriver Storage

// RegisterDriver ثبت Driver جدید در Registry
func RegisterDriver(name string, driver Storage) {
	Registry[name] = driver
}

// SelectDriver انتخاب Driver بر اساس نام
func SelectDriver(name string) (Storage, error) {
	if name == "" {
		if DefaultDriver == nil {
			return nil, fmt.Errorf("no default storage driver set")
		}
		return DefaultDriver, nil
	}
	driver, exists := Registry[name]
	if !exists {
		return nil, fmt.Errorf("driver not found: %s", name)
	}
	return driver, nil
}

// SetDefaultDriver تنظیم درایور پیش‌فرض
func SetDefaultDriver(name string) error {
	driver, exists := Registry[name]
	if !exists {
		return fmt.Errorf("driver not found: %s", name)
	}
	DefaultDriverName = name
	DefaultDriver = driver
	return nil
}
