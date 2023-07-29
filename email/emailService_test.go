package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmailService(t *testing.T) {
	t.Run("with valid config", func(t *testing.T) {
		// Assuming a valid config.yaml file exists in the specified path
		es, err := NewEmailService()
		assert.NoError(t, err)
		assert.NotNil(t, es)
	})
	t.Run("with invalid config", func(t *testing.T) {
		// Assuming an invalid config.yaml file exists in the specified path
		es, err := NewEmailService()
		assert.Error(t, err)
		assert.Nil(t, es)
	})
}
func TestSetFrom(t *testing.T) {
	t.Run("with valid email", func(t *testing.T) {
		es := &EmailService{}
		es.SetFrom("test@example.com")
		assert.Equal(t, "test@example.com", es.from)
	})
	t.Run("with empty email", func(t *testing.T) {
		es := &EmailService{}
		es.SetFrom("")
		assert.Equal(t, "", es.from)
	})
}
