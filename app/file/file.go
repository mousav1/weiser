package file

import (
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
)

// WriteToFile writes data to a file
func WriteToFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

// ReadFromFile reads data from a file
func ReadFromFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// DeleteFile deletes a file
func DeleteFile(filename string) error {
	return os.Remove(filename)
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
		return 0, err
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
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		err = c.SaveFile(file, uploadDir+"/"+file.Filename)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
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
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		defer file.Close()

		// Send the file to the client
		return c.SendStream(file)
	}
}
