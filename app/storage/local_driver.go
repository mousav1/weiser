package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type LocalDriver struct {
	BasePath string
}

func NewLocalDriver(basePath string) *LocalDriver {
	return &LocalDriver{BasePath: basePath}
}

func (ld *LocalDriver) Put(path string, content io.ReadSeeker) error {
	filePath := filepath.Join(ld.BasePath, path)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, content)
	return err
}

func (ld *LocalDriver) Get(path string) (io.Reader, error) {
	filePath := filepath.Join(ld.BasePath, path)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (ld *LocalDriver) Delete(path string) error {
	filePath := filepath.Join(ld.BasePath, path)
	return os.Remove(filePath)
}

func (ld *LocalDriver) Exists(path string) (bool, error) {
	filePath := filepath.Join(ld.BasePath, path)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func (ld *LocalDriver) List(directory string) ([]string, error) {
	dirPath := filepath.Join(ld.BasePath, directory)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func (ld *LocalDriver) Missing(path string) (bool, error) {
	exists, err := ld.Exists(path)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

func (ld *LocalDriver) Download(path string) (io.Reader, error) {
	return ld.Get(path)
}

func (ld *LocalDriver) URL(path string) (string, error) {
	return fmt.Sprintf("%s/%s", ld.BasePath, path), nil
}

func (ld *LocalDriver) TemporaryURL(path string, expiresIn int64) (string, error) {
	filePath := filepath.Join(ld.BasePath, path)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	expireTime := time.Now().Add(time.Duration(expiresIn) * time.Second)
	signature := ld.generateSignature(path, expireTime, fileInfo.ModTime())

	signedURL, err := url.Parse(fmt.Sprintf("%s/%s", ld.BasePath, path))
	if err != nil {
		return "", err
	}

	query := signedURL.Query()
	query.Set("expires", expireTime.Format(time.RFC3339))
	query.Set("signature", signature)
	signedURL.RawQuery = query.Encode()

	return signedURL.String(), nil
}

func (ld *LocalDriver) generateSignature(path string, expireTime time.Time, modTime time.Time) string {
	// در اینجا کلید امنیتی مشترک را فرض می کنیم
	secretKey := []byte("your_secret_key")

	// ساخت امضا با استفاده از HMAC-SHA256
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(path))
	h.Write([]byte(expireTime.Format(time.RFC3339)))
	h.Write([]byte(modTime.Format(time.RFC3339)))
	signature := hex.EncodeToString(h.Sum(nil))

	return signature
}
func (ld *LocalDriver) Size(path string) (int64, error) {
	filePath := filepath.Join(ld.BasePath, path)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func (ld *LocalDriver) Copy(sourcePath string, destinationPath string) error {
	srcPath := filepath.Join(ld.BasePath, sourcePath)
	dstPath := filepath.Join(ld.BasePath, destinationPath)
	input, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer output.Close()
	_, err = io.Copy(output, input)
	return err
}

func (ld *LocalDriver) Move(sourcePath string, destinationPath string) error {
	err := ld.Copy(sourcePath, destinationPath)
	if err != nil {
		return err
	}
	return ld.Delete(sourcePath)
}
