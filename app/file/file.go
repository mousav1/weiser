package file

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
)

// WriteToFile writes data to a file
func WriteToFile(filename string, data []byte) error {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return errors.New("failed to write data to file")
	}
	return nil
}

// ReadFromFile reads data from a file
func ReadFromFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("failed to read data from file")
	}
	return data, nil
}

// DeleteFile deletes a file
func DeleteFile(filename string) error {
	if err := os.Remove(filename); err != nil {
		return errors.New("failed to delete file")
	}
	return nil
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(filename string) (int64, error) {
	file, err := os.Stat(filename)
	if err != nil {
		return 0, errors.New("failed to get file size")
	}
	return file.Size(), nil
}

// UploadFile returns a handler function for uploading files
func UploadFile(uploadDir string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check request method
		if c.Method() != "POST" {
			return c.Status(fiber.StatusMethodNotAllowed).SendString("Method not allowed")
		}

		// Save file to the specified directory
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("failed to upload file")
		}
		err = c.SaveFile(file, uploadDir+"/"+file.Filename)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("failed to save file")
		}

		// Return the saved file name
		return c.SendString(file.Filename)
	}
}

// DownloadFile returns a handler function for downloading files
func DownloadFile(uploadDir string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Open the file
		filename := c.Params("filename")
		file, err := os.Open(uploadDir + "/" + filename)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("file not found")
		}
		defer file.Close()

		// Send the file to the client
		return c.SendStream(file)
	}
}
