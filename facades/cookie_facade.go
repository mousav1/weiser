package facades

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mousav1/weiser/app/cookies"
)

// CookieFacade provides a simplified interface to interact with cookies.
type CookieFacade struct{}

// Set sets a cookie with the given name, value, and expiration time.
func (cf *CookieFacade) Set(c *fiber.Ctx, name string, value string, expire time.Time) {
	cookies.SetCookie(c, name, value, expire)
}

// Get retrieves a cookie by name.
func (cf *CookieFacade) Get(c *fiber.Ctx, name string) (string, error) {
	return cookies.GetCookie(c, name)
}

func NewCookieFacade() *CookieFacade {
	return &CookieFacade{}
}
