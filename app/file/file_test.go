package file

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestWriteToFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal("failed to create temporary file:", err)
	}
	defer os.Remove(tempFile.Name())

	// Test WriteToFile function
	data := []byte("test data")
	err = WriteToFile(tempFile.Name(), data)
	if err != nil {
		t.Errorf("WriteToFile returned an error: %v", err)
	}

	// Read the file to verify the written data
	fileData, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal("failed to read file:", err)
	}

	// Check if the read data matches the written data
	if string(fileData) != string(data) {
		t.Errorf("data mismatch, expected: %s, got: %s", string(data), string(fileData))
	}
}

func TestReadFromFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal("failed to create temporary file:", err)
	}
	defer os.Remove(tempFile.Name())

	// Write data to the file
	data := []byte("test data")
	err = ioutil.WriteFile(tempFile.Name(), data, 0644)
	if err != nil {
		t.Fatal("failed to write data to file:", err)
	}

	// Test ReadFromFile function
	fileData, err := ReadFromFile(tempFile.Name())
	if err != nil {
		t.Errorf("ReadFromFile returned an error: %v", err)
	}

	// Check if the read data matches the expected data
	if string(fileData) != string(data) {
		t.Errorf("data mismatch, expected: %s, got: %s", string(data), string(fileData))
	}
}

func TestDeleteFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal("failed to create temporary file:", err)
	}
	defer os.Remove(tempFile.Name())

	// Test DeleteFile function
	err = DeleteFile(tempFile.Name())
	if err != nil {
		t.Errorf("DeleteFile returned an error: %v", err)
	}

	// Check if the file still exists after deletion
	if FileExists(tempFile.Name()) {
		t.Errorf("file still exists after deletion")
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal("failed to create temporary file:", err)
	}
	defer os.Remove(tempFile.Name())

	// Test FileExists function
	exists := FileExists(tempFile.Name())
	if !exists {
		t.Errorf("FileExists returned false for an existing file")
	}
}

func TestGetFileSize(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal("failed to create temporary file:", err)
	}
	defer os.Remove(tempFile.Name())

	// Write data to the file
	data := []byte("test data")
	err = ioutil.WriteFile(tempFile.Name(), data, 0644)
	if err != nil {
		t.Fatal("failed to write data to file:", err)
	}

	// Test GetFileSize function
	fileSize, err := GetFileSize(tempFile.Name())
	if err != nil {
		t.Errorf("GetFileSize returned an error: %v", err)
	}

	// Check if the file size matches the expected size
	expectedSize := int64(len(data))
	if fileSize != expectedSize {
		t.Errorf("file size mismatch, expected: %d, got: %d", expectedSize, fileSize)
	}
}
func TestUploadFile(t *testing.T) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create a temporary directory for testing
	tempDir := filepath.Join(currentDir, "uploads")
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	fileContent := []byte("test content")
	filePath := filepath.Join(tempDir, "test.txt")
	err = ioutil.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filePath)

	// Create a mock Fiber context
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Set the request method to POST
	ctx.Method("POST")

	// Set the content type in the response
	ctx.Type("multipart/form-data")

	// Set the form file in the request
	ctx.Request().Header.Set("Content-Disposition", `form-data; name="file"; filename="`+filePath+`"`)
	ctx.Request().SetBody([]byte("test content"))

	// Call the UploadFile handler
	err = UploadFile(tempDir)(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the file was saved correctly
	uploadedFile := filepath.Join(tempDir, "test.txt")
	_, err = os.Stat(uploadedFile)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the saved file has the correct content
	content, err := ioutil.ReadFile(uploadedFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "test content" {
		t.Errorf("expected file content 'test content', got '%s'", string(content))
	}
}

func TestDownloadFile(t *testing.T) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create a temporary directory for testing
	tempDir := filepath.Join(currentDir, "uploads")
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	fileContent := []byte("test content")
	filePath := filepath.Join(tempDir, "test.txt")
	err = ioutil.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filePath)

	// Create a Fiber app
	app := fiber.New()

	// Register the DownloadFile handler
	app.Get("/files/:filename", DownloadFile(tempDir))

	// Create a test request
	req := &fasthttp.Request{}
	req.SetRequestURI("/files/test.txt")

	// Create a test request context
	ctx := &fasthttp.RequestCtx{}
	ctx.Init(req, nil, nil)

	// Process the request
	app.Handler()(ctx)

	// Check the response status code
	if ctx.Response.StatusCode() != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, ctx.Response.StatusCode())
	}

	// Check the response body
	expectedBody, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(ctx.Response.Body()) != string(expectedBody) {
		t.Errorf("expected response body %s, got %s", string(expectedBody), string(ctx.Response.Body()))
	}
}
