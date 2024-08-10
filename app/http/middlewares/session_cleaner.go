package middleware

import (
	"log"
	"time"

	"github.com/mousav1/weiser/app/session"
)

func StartSessionCleaner() {
	sessionManager := session.GetSessionManager()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := cleanExpiredSessions(sessionManager); err != nil {
			log.Printf("Error cleaning expired sessions: %v", err)
		}
	}
}

func cleanExpiredSessions(manager *session.SessionManager) error {
	sessionIDs, err := manager.GetSessionIDs()
	if err != nil {
		return err
	}

	for _, sessionID := range sessionIDs {
		if err := manager.CheckExpiration(sessionID); err != nil {
			if clearErr := manager.Clear(sessionID); clearErr != nil {
				log.Printf("Error deleting session %v: %v", sessionID, clearErr)
			} else {
				log.Printf("Deleted session with ID %v", sessionID)
			}
		}
	}

	return nil
}
