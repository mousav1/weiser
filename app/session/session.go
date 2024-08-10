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
	"github.com/mousav1/weiser/database"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type SessionManager struct {
	storage Storage
}

type Session struct {
	ID      string
	Data    map[string]interface{}
	Expires time.Time
}

var sm *SessionManager

// InitSessionManager initializes the session manager with the specified storage type and configuration
func InitSessionManager() error {
	storageType := viper.GetString("session.type")
	var storage Storage
	var err error

	switch storageType {
	case "redis":
		storage, err = NewRedisStorage()
	case "file":
		filePath := viper.GetString("session.file.path")
		storage, err = NewFileStorage(filePath)
	case "memory":
		storage = &InMemoryStorage{sessions: make(map[string]Session)}
	default:
		return errors.New("invalid session storage type")
	}

	if err != nil {
		return err
	}

	sm = NewSessionManager(storage)
	return nil
}

// NewRedisStorage initializes Redis storage
func NewRedisStorage() (*RedisStorage, error) {
	redisConfig := viper.GetStringMapString("database.redis.cache")
	client, err := database.ConnectToRedis(redisConfig)
	if err != nil {
		return nil, err
	}
	return &RedisStorage{client: client}, nil
}

// GetSessionManager returns the initialized session manager
func GetSessionManager() *SessionManager {
	return sm
}

// NewSessionManager creates a new session manager with the given storage
func NewSessionManager(storage Storage) *SessionManager {
	return &SessionManager{storage: storage}
}

// StartSession creates a new session and sets a cookie for it
func (sm *SessionManager) StartSession(c *fiber.Ctx) *Session {
	sessionID, err := generateSessionID()
	if err != nil {
		log.Println("Failed to generate session ID:", err)
		return nil
	}

	expirationTime := viper.GetDuration("session.expirationTime")
	expiration := time.Now().Add(expirationTime)

	sessionData := Session{
		ID:      sessionID,
		Data:    make(map[string]interface{}),
		Expires: expiration,
	}

	if err := sm.storage.Set(sessionID, sessionData); err != nil {
		log.Println("Failed to store session:", err)
		return nil
	}

	cookies.SetCookie(c, "weiser_session", sessionID, expiration)
	return &sessionData
}

// Set updates a session value
func (sm *SessionManager) Set(key string, value interface{}, sessionID string) (*Session, error) {
	if err := sm.CheckExpiration(sessionID); err != nil {
		return nil, err
	}

	sessionData, err := sm.GetDataBySessionID(sessionID)
	if err != nil {
		return nil, err
	}
	sessionData.Data[key] = value
	sessionData.Expires = time.Now().Add(viper.GetDuration("session.expirationTime"))

	if err := sm.storage.Set(sessionID, *sessionData); err != nil {
		return nil, err
	}

	return sessionData, nil
}

// Get retrieves a session value
func (sm *SessionManager) Get(key string, sessionID string) interface{} {
	if err := sm.CheckExpiration(sessionID); err != nil {
		return nil
	}

	sessionData, err := sm.GetDataBySessionID(sessionID)
	if err != nil {
		return nil
	}
	return sessionData.Data[key]
}

// GetDataBySessionID retrieves session data by its ID
func (sm *SessionManager) GetDataBySessionID(sessionID string) (*Session, error) {
	sessionData, err := sm.storage.Get(sessionID)
	if err != nil {
		return nil, err
	}
	return &sessionData, nil
}

// Delete removes a session value
func (sm *SessionManager) Delete(key string, sessionID string) {
	if err := sm.CheckExpiration(sessionID); err != nil {
		return
	}

	sessionData, err := sm.GetDataBySessionID(sessionID)
	if err != nil {
		return
	}

	delete(sessionData.Data, key)
	sessionData.Expires = time.Now().Add(viper.GetDuration("session.expirationTime"))

	if err := sm.storage.Set(sessionID, *sessionData); err != nil {
		log.Println("Failed to update session:", err)
	}
}

// Clear removes a session
func (sm *SessionManager) Clear(sessionID string) error {
	return sm.storage.Delete(sessionID)
}

// CheckExpiration checks if the session has expired
func (sm *SessionManager) CheckExpiration(sessionID string) error {
	sessionData, err := sm.GetDataBySessionID(sessionID)
	if err != nil {
		return err
	}

	if time.Now().After(sessionData.Expires) {
		if err := sm.Clear(sessionID); err != nil {
			return err
		}
		return errors.New("session has expired")
	}

	return nil
}

// UpdateExpiration updates the expiration time of a session
func (sm *SessionManager) UpdateExpiration(sessionID string, expiration time.Duration) error {
	sessionData, err := sm.GetDataBySessionID(sessionID)
	if err != nil {
		return err
	}

	sessionData.Expires = time.Now().Add(expiration)

	return sm.storage.Set(sessionID, *sessionData)
}

// IsValid checks if a session is valid
func (s *Session) IsValid() bool {
	return time.Now().Before(s.Expires)
}

// GetSessionIDs retrieves all session IDs
func (sm *SessionManager) GetSessionIDs() ([]string, error) {
	sessionIDs, err := sm.storage.GetSessionIDs()
	if err != nil {
		return []string{}, nil
	}
	return sessionIDs, nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Storage interface for session storage operations
type Storage interface {
	Set(key string, value Session) error
	Get(key string) (Session, error)
	Delete(key string) error
	GetSessionIDs() ([]string, error)
}

// InMemoryStorage implementation for in-memory session storage
type InMemoryStorage struct {
	sessions map[string]Session
}

func (ims *InMemoryStorage) Set(key string, value Session) error {
	ims.sessions[key] = value
	return nil
}

func (ims *InMemoryStorage) Get(key string) (Session, error) {
	session, ok := ims.sessions[key]
	if !ok {
		return Session{}, errors.New("session not found")
	}
	return session, nil
}

func (ims *InMemoryStorage) Delete(key string) error {
	delete(ims.sessions, key)
	return nil
}

func (ims *InMemoryStorage) GetSessionIDs() ([]string, error) {
	sessionIDs := make([]string, 0, len(ims.sessions))
	for sessionID := range ims.sessions {
		sessionIDs = append(sessionIDs, sessionID)
	}
	return sessionIDs, nil
}

// RedisStorage implementation for Redis session storage
type RedisStorage struct {
	client *redis.Client
}

func (rs *RedisStorage) Set(key string, value Session) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rs.client.Set(context.Background(), key, jsonValue, 0).Err()
}

func (rs *RedisStorage) Get(key string) (Session, error) {
	value, err := rs.client.Get(context.Background(), key).Result()
	if err != nil {
		return Session{}, err
	}
	var sessionData Session
	err = json.Unmarshal([]byte(value), &sessionData)
	if err != nil {
		return Session{}, err
	}
	return sessionData, nil
}

func (rs *RedisStorage) Delete(key string) error {
	return rs.client.Del(context.Background(), key).Err()
}

func (rs *RedisStorage) GetSessionIDs() ([]string, error) {
	sessionIDs, err := rs.client.Keys(context.Background(), "*").Result()
	if err != nil {
		return []string{}, err
	}
	return sessionIDs, nil
}

// FileStorage implementation for file-based session storage
type FileStorage struct {
	filePath string
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return nil, err
		}
	}

	return &FileStorage{filePath: filePath}, nil
}

func (fs *FileStorage) Set(key string, value Session) error {
	sessions, err := fs.readSessions()
	if err != nil {
		return err
	}

	sessions[key] = value
	return fs.writeSessions(sessions)
}

func (fs *FileStorage) Get(key string) (Session, error) {
	sessions, err := fs.readSessions()
	if err != nil {
		return Session{}, err
	}

	sessionData, ok := sessions[key]
	if !ok {
		return Session{}, errors.New("session not found")
	}

	return sessionData, nil
}

func (fs *FileStorage) Delete(key string) error {
	sessions, err := fs.readSessions()
	if err != nil {
		return err
	}

	delete(sessions, key)
	return fs.writeSessions(sessions)
}

func (fs *FileStorage) GetSessionIDs() ([]string, error) {
	sessions, err := fs.readSessions()
	if err != nil {
		return nil, err
	}

	sessionIDs := make([]string, 0, len(sessions))
	for sessionID := range sessions {
		sessionIDs = append(sessionIDs, sessionID)
	}
	return sessionIDs, nil
}

func (fs *FileStorage) readSessions() (map[string]Session, error) {
	defaultSessions := map[string]Session{"empty": {}}
	defaultBytes, _ := json.Marshal(defaultSessions)

	fi, err := os.Stat(fs.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(fs.filePath)
			if err != nil {
				return nil, err
			}
			defer file.Close()

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
	defer file.Close()

	decoder := json.NewDecoder(file)
	sessions := make(map[string]Session)
	if err = decoder.Decode(&sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (fs *FileStorage) writeSessions(sessions map[string]Session) error {
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
