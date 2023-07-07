package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/cookies"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type SessionManager struct {
	storage Storage
}

type Session struct {
	ID   string
	Data map[string]interface{}
}

var sm *SessionManager

// InitSessionManager initializes the session manager with the specified storage type and configuration
func InitSessionManager() error {
	storageType := viper.GetString("session.type")

	var storage Storage
	var err error

	switch storageType {
	case "redis":
		var err error
		storage, err = NewRedisStorage()
		if err != nil {
			return err
		}
	case "file":
		filePath := viper.GetString("session.file.path")
		storage, err = NewFileStorage(filePath)
		if err != nil {
			return err
		}
	case "memory":
		storage = &InMemoryStorage{sessions: make(map[string]map[string]interface{})}
	default:
		return errors.New("invalid session storage type")
	}

	sm = NewSessionManager(storage)

	return nil
}

func NewRedisStorage() (*RedisStorage, error) {
	redisAddr := viper.GetString("redis.addr")
	redisPassword := viper.GetString("redis.password")
	redisDB := viper.GetInt("redis.db")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test the Redis connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisStorage{client: client}, nil
}

// GetSessionManager returns the initialized session manager
func GetSessionManager() *SessionManager {
	return sm
}

func NewSessionManager(storage Storage) *SessionManager {
	return &SessionManager{storage: storage}
}

func (sm *SessionManager) StartSession(c *fiber.Ctx) *Session {
	sessionID, err := generateSessionID()
	if err != nil {
		log.Println("Failed to generate session ID:", err)
		return nil
	}

	sessionData := make(map[string]interface{})

	// Set the expiration time for the session
	expirationTime := viper.GetDuration("session.expirationTime")
	expiration := time.Now().Add(expirationTime)
	sessionData["expirationTime"] = expiration.Format(time.RFC3339)

	sm.storage.Set(sessionID, sessionData)

	// Create a new cookie with the session ID
	cookies.SetCookie(c, "weiser_session", sessionID, expiration)

	return &Session{ID: sessionID, Data: sessionData}
}

func (sm *SessionManager) Set(key string, value interface{}, sessionID string) (*Session, error) {
	sm.CheckExpiration(sessionID)

	sessionData, err := sm.GetDataBySessionID(sessionID)
	if err != nil {
		return nil, err
	}
	sessionData[key] = value
	sm.storage.Set(sessionID, sessionData)
	return &Session{ID: sessionID, Data: sessionData}, nil
}

func (sm *SessionManager) Get(key string, sessionID string) interface{} {
	sm.CheckExpiration(sessionID)

	sessionData, err := sm.GetDataBySessionID(sessionID)
	if err != nil {
		return nil
	}
	value, _ := sessionData[key]
	return value
}

func (sm *SessionManager) GetDataBySessionID(sessionID interface{}) (map[string]interface{}, error) {
	sessionData, err := sm.storage.Get(sessionID.(string))
	if err != nil {
		return nil, err
	}
	return sessionData, nil
}

func (sm *SessionManager) Delete(key string, sessionID string) {
	sm.CheckExpiration(sessionID)

	sessionData, err := sm.GetDataBySessionID(sessionID)
	if err != nil {
		return
	}
	delete(sessionData, key)
	sm.storage.Set(sessionID, sessionData)
}

func (sm *SessionManager) Clear(sessionID string) error {
	err := sm.storage.Delete(sessionID)
	if err != nil {
		return err
	}
	return nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type Storage interface {
	Set(key string, value map[string]interface{}) error
	Get(key string) (map[string]interface{}, error)
	Delete(key string) error
	GetSessionIDs() ([]string, error)
}

type InMemoryStorage struct {
	sessions map[string]map[string]interface{}
}

func (ims *InMemoryStorage) Set(key string, value map[string]interface{}) error {
	ims.sessions[key] = value
	return nil
}

func (ims *InMemoryStorage) Get(key string) (map[string]interface{}, error) {
	value, ok := ims.sessions[key]
	if !ok {
		return nil, errors.New("session not found")
	}
	return value, nil
}

func (ims *InMemoryStorage) Delete(key string) error {
	delete(ims.sessions, key)
	return nil
}

type RedisStorage struct {
	client *redis.Client
}

func (rs *RedisStorage) Set(key string, value map[string]interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rs.client.Set(context.Background(), key, jsonValue, 0).Err()
}

func (rs *RedisStorage) Get(key string) (map[string]interface{}, error) {
	value, err := rs.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	var sessionData map[string]interface{}
	err = json.Unmarshal([]byte(value), &sessionData)
	if err != nil {
		return nil, err
	}
	return sessionData, nil
}

func (rs *RedisStorage) Delete(key string) error {
	return rs.client.Del(context.Background(), key).Err()
}

func (m *SessionManager) CheckExpiration(sessionID string) error {
	// Get the session data from the storage using the session ID
	data, err := m.GetDataBySessionID(sessionID)
	if err != nil {
		return err
	}

	// Get the session expiration time
	expirationTimeString, ok := data["expirationTime"].(string)
	if !ok {
		return errors.New("invalid expiration time")
	}
	expirationTime, err := time.Parse(time.RFC3339, expirationTimeString)
	if err != nil {
		return err
	}

	// Check if the session has expired
	if time.Now().After(expirationTime) {
		// If the session has expired, delete it from the storage
		if err = m.Clear(sessionID); err != nil {
			return err
		}
		return errors.New("session has expired")
	}

	// If the session has not expired, update its expiration time
	if err = m.UpdateExpiration(sessionID, viper.GetDuration("session.expirationTime")); err != nil {
		return err
	}

	return nil
}

func (sm *SessionManager) UpdateExpiration(sessionID string, expiration time.Duration) error {
	// Get the session data from the storage using the session ID
	sessionData, err := sm.GetDataBySessionID(sessionID)
	if err != nil {
		return err
	}

	// Set the new expiration time for the session data
	expirationTime := time.Now().Add(expiration)
	sessionData["expirationTime"] = expirationTime

	// Update the session data in the storage
	if err := sm.storage.Set(sessionID, sessionData); err != nil {
		return err
	}

	return nil
}

func (s *Session) IsValid() bool {
	if s == nil {
		return false
	}

	expirationTime, ok := s.Data["expirationTime"].(time.Time)
	if !ok {
		return false
	}

	return time.Now().Before(expirationTime)
}

func (sm *SessionManager) GetSessionIDs() []string {
	sessionIDs, err := sm.storage.GetSessionIDs()
	if err != nil {
		return []string{}
	}
	return sessionIDs
}

func (ims *InMemoryStorage) GetSessionIDs() ([]string, error) {
	sessionIDs := make([]string, 0, len(ims.sessions))
	for sessionID := range ims.sessions {
		sessionIDs = append(sessionIDs, sessionID)
	}
	return sessionIDs, nil
}

func (rs *RedisStorage) GetSessionIDs() ([]string, error) {
	sessionIDs, err := rs.client.Keys(context.Background(), "*").Result()
	if err != nil {
		return []string{}, err
	}
	return sessionIDs, nil
}

type FileStorage struct {
	filePath string
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	// Check if the file exists or not
	_, err := os.Stat(filePath)
	if err != nil {
		// If the file does not exist, create it
		if os.IsNotExist(err) {
			_, err = os.Create(filePath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &FileStorage{filePath: filePath}, nil
}

func (fs *FileStorage) Set(key string, value map[string]interface{}) error {
	// Read the existing sessions from the file
	sessions, err := fs.readSessions()
	if err != nil {
		return err
	}

	// If the key does not exist, create a new session
	if _, ok := sessions[key]; !ok {
		sessions[key] = make(map[string]interface{})
	}

	// Update the sessions with the new data
	for k, v := range value {
		sessions[key][k] = v
	}

	// Write the updated sessions to the file
	return fs.writeSessions(sessions)
}

func (fs *FileStorage) Get(key string) (map[string]interface{}, error) {
	// Read the existing sessions from the file
	sessions, err := fs.readSessions()
	if err != nil {
		return nil, err
	}

	// Get the session data
	sessionData, ok := sessions[key]
	if !ok {
		return nil, errors.New("session not found")
	}

	return sessionData, nil
}

func (fs *FileStorage) Delete(key string) error {
	// Read the existing sessions from the file
	sessions, err := fs.readSessions()
	if err != nil {
		return err
	}

	// Delete the session data
	delete(sessions, key)

	// Write the updated sessions to the file
	return fs.writeSessions(sessions)
}

func (fs *FileStorage) GetSessionIDs() ([]string, error) {
	// Read the existing sessions from the file
	sessions, err := fs.readSessions()
	if err != nil {
		return nil, err
	}

	// Get the session IDs
	sessionIDs := make([]string, 0, len(sessions))
	for sessionID := range sessions {
		sessionIDs = append(sessionIDs, sessionID)
	}

	return sessionIDs, nil
}

func (fs *FileStorage) readSessions() (map[string]map[string]interface{}, error) {

	// Define the default sessions
	defaultSessions := map[string]map[string]interface{}{"empty": nil}
	defaultBytes, _ := json.Marshal(defaultSessions)

	fi, err := os.Stat(fs.filePath)
	if err != nil {
		// If the file does not exist, create it
		if os.IsNotExist(err) {
			file, err := os.Create(fs.filePath)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			// Write the default sessions to the file
			if _, err := file.Write(defaultBytes); err != nil {
				return nil, err
			}

			return defaultSessions, nil
		} else {
			return nil, err
		}
	}
	if fi.Size() == 0 {
		return defaultSessions, nil
	}

	file, err := os.OpenFile(fs.filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Close the file after reading the data

	decoder := json.NewDecoder(file)
	sessions := make(map[string]map[string]interface{})
	if err = decoder.Decode(&sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (fs *FileStorage) writeSessions(sessions map[string]map[string]interface{}) error {
	// Write the updated sessions to the file
	file, err := os.OpenFile(fs.filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(sessions); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	return nil
}
