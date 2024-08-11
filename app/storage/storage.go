package storage

import (
	"io"
)

// Storage interface برای مدیریت فایل‌ها
type Storage interface {
	Put(path string, content io.ReadSeeker) error
	Get(path string) (io.Reader, error)
	Delete(path string) error
	Exists(path string) (bool, error)
	List(directory string) ([]string, error)
	Missing(path string) (bool, error)
	Download(path string) (io.Reader, error)
	URL(path string) (string, error)
	TemporaryURL(path string, expiresIn int64) (string, error)
	Size(path string) (int64, error)
	Copy(sourcePath string, destinationPath string) error
	Move(sourcePath string, destinationPath string) error
}
