package session

import (
	"errors"
	"log"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func BeforeEach() {
	// Perform initial setup for tests
	viper.SetConfigFile("../../config/config.yaml")
	viper.Set("session.file.path", "../../storage/logs/logs.txt")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read configuration file: %s", err)
	}
}

// MockStorage implements the Storage interface for testing purposes
type MockStorage struct {
	data map[string]map[string]interface{}
}

func (ms *MockStorage) Set(key string, value map[string]interface{}) error {
	ms.data[key] = value
	return nil
}

func (ms *MockStorage) Get(key string) (map[string]interface{}, error) {
	value, ok := ms.data[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (ms *MockStorage) Delete(key string) error {
	delete(ms.data, key)
	return nil
}

func (ms *MockStorage) GetSessionIDs() ([]string, error) {
	var sessionIDs []string
	for sessionID := range ms.data {
		sessionIDs = append(sessionIDs, sessionID)
	}
	return sessionIDs, nil
}

func TestSessionManager_StartSession(t *testing.T) {
	BeforeEach()

	// Initialize the session manager with the mock storage
	storage := &MockStorage{data: make(map[string]map[string]interface{})}
	sm := NewSessionManager(storage)

	// Create a mock Fiber context
	ctx := fiber.New().AcquireCtx(&fasthttp.RequestCtx{})

	// Start a session
	session := sm.StartSession(ctx)

	// Check if the session is created successfully
	if session == nil {
		t.Error("Failed to start session")
	}

	// Check if the session ID is generated and set correctly
	if session.ID == "" {
		t.Error("Invalid session ID")
	}

	// Check if the session data is stored in the storage correctly
	sessionData, err := storage.Get(session.ID)
	if err != nil {
		t.Errorf("Failed to get session data: %v", err)
	}
	if sessionData == nil {
		t.Error("Session data not stored")
	}
}

func TestSessionManager_Set_Get_Delete(t *testing.T) {
	BeforeEach()

	// Initialize the session manager with the mock storage
	storage := &MockStorage{data: make(map[string]map[string]interface{})}
	sm := NewSessionManager(storage)

	// Create a mock session
	session := &Session{
		ID: "session_id",
		Data: map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}

	storage.data[session.ID] = session.Data

	// Set a value in the session
	_, err := sm.Set("key3", "value3", session.ID)
	if err != nil {
		t.Errorf("Failed to set value in session: %v", err)
	}

	// Get a value from the session
	value := sm.Get("key1", session.ID)
	if value != "value1" {
		t.Error("Failed to get value from session")
	}

	// Delete a value from the session
	sm.Delete("key2", session.ID)

	// Check if the value is deleted from the session
	value = sm.Get("key2", session.ID)
	if value != nil {
		t.Error("Failed to delete value from session")
	}
}

// You can write more tests for other methods in the SessionManager struct

func TestGenerateSessionID(t *testing.T) {
	BeforeEach()

	// Generate a session ID
	sessionID, err := generateSessionID()
	if err != nil {
		t.Errorf("Failed to generate session ID: %v", err)
	}

	// Check if the generated session ID is valid
	if sessionID == "" {
		t.Error("Invalid session ID")
	}
}

func TestInitSessionManager(t *testing.T) {
	BeforeEach()

	// Initialize the session manager
	err := InitSessionManager()
	if err != nil {
		t.Errorf("Failed to initialize session manager: %v", err)
	}

	// Get the session manager
	sm := GetSessionManager()

	// Check if the session manager is initialized
	if sm == nil {
		t.Error("Session manager not initialized")
	}
}
