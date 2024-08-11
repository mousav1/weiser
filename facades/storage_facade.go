package facades

import (
	"io"

	"github.com/mousav1/weiser/app/storage"
)

// StorageFacade نماینده Facade برای سیستم ذخیره‌سازی
type StorageFacade struct{}

// NewStorageFacade ایجاد یک نمونه جدید از StorageFacade
func NewStorageFacade() *StorageFacade {
	return &StorageFacade{}
}

// GetDefaultDriver دریافت درایور پیش‌فرض
func (sf *StorageFacade) GetDefaultDriver() storage.Storage {
	return storage.DefaultDriver
}

// Put فایل را در درایور پیش‌فرض قرار می‌دهد
func (sf *StorageFacade) Put(path string, content io.ReadSeeker) error {
	return storage.DefaultDriver.Put(path, content)
}

// Get فایل را از درایور پیش‌فرض دریافت می‌کند
func (sf *StorageFacade) Get(path string) (io.Reader, error) {
	return storage.DefaultDriver.Get(path)
}

// Delete فایل را از درایور پیش‌فرض حذف می‌کند
func (sf *StorageFacade) Delete(path string) error {
	return storage.DefaultDriver.Delete(path)
}

// Exists بررسی می‌کند که آیا فایل در درایور پیش‌فرض وجود دارد یا خیر
func (sf *StorageFacade) Exists(path string) (bool, error) {
	return storage.DefaultDriver.Exists(path)
}

// List فایل‌های موجود در دایرکتوری در درایور پیش‌فرض را لیست می‌کند
func (sf *StorageFacade) List(directory string) ([]string, error) {
	return storage.DefaultDriver.List(directory)
}

// URL ایجاد URL دائمی برای دسترسی به فایل در درایور پیش‌فرض
func (sf *StorageFacade) URL(path string) (string, error) {
	return storage.DefaultDriver.URL(path)
}

// TemporaryURL ایجاد URL موقتی برای دسترسی به فایل در درایور پیش‌فرض
func (sf *StorageFacade) TemporaryURL(path string, expiresIn int64) (string, error) {
	return storage.DefaultDriver.TemporaryURL(path, expiresIn)
}

// Size اندازه فایل را در درایور پیش‌فرض برمی‌گرداند
func (sf *StorageFacade) Size(path string) (int64, error) {
	return storage.DefaultDriver.Size(path)
}

// Copy فایل را از مسیر مبدا به مسیر مقصد در درایور پیش‌فرض کپی می‌کند
func (sf *StorageFacade) Copy(sourcePath string, destinationPath string) error {
	return storage.DefaultDriver.Copy(sourcePath, destinationPath)
}

// Move فایل را از مسیر مبدا به مسیر مقصد در درایور پیش‌فرض منتقل می‌کند
func (sf *StorageFacade) Move(sourcePath string, destinationPath string) error {
	return storage.DefaultDriver.Move(sourcePath, destinationPath)
}
