package cookies

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type cookieConfig struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	Expires  time.Time
	Secure   bool
	SameSite http.SameSite
	HTTPOnly bool
}

var config cookieConfig

func init() {
	cookie := viper.Get("cookie")
	if cookie != nil {
		cookieMap := cookie.(map[string]interface{})
		config.Name = getStringOrDefault(cookieMap["name"], "cookie_name")
		config.Value = ""
		config.Path = getStringOrDefault(cookieMap["path"], "/")
		config.Domain = getStringOrDefault(cookieMap["domain"], "")
		config.Expires, _ = time.Parse(time.RFC3339, getStringOrDefault(cookieMap["expires"], "2030-12-31T00:00:00Z"))
		config.Secure = getBoolOrDefault(cookieMap["secure"], true)
		sameSite, _ := strconv.Atoi(getStringOrDefault(cookieMap["samesite"], "0"))
		config.SameSite = http.SameSite(sameSite)
		config.HTTPOnly = getBoolOrDefault(cookieMap["httponly"], true)
	}
}

func SetCookie(c *fiber.Ctx, name string, value string, expire time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     config.Path,
		Domain:   config.Domain,
		Expires:  expire,
		Secure:   config.Secure,
		SameSite: sameSiteToString(config.SameSite),
		HTTPOnly: config.HTTPOnly,
	})
}

func GetCookie(c *fiber.Ctx, name string) (interface{}, error) {
	cookie := c.Cookies(name)
	if cookie == "" {
		return nil, errors.New("cookie not found")
	}
	return cookie, nil
}

func SetHttpCookie(w http.ResponseWriter, name string, value interface{}, expire time.Duration, secure bool, httpOnly bool, sameSite http.SameSite) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    fmt.Sprint(value),
		Path:     config.Path,
		Domain:   config.Domain,
		Expires:  time.Now().Add(expire),
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
	})
}

func sameSiteToString(sameSite http.SameSite) string {
	switch sameSite {
	case http.SameSiteDefaultMode:
		return "Lax"
	case http.SameSiteStrictMode:
		return "Strict"
	case http.SameSiteLaxMode:
		return "Lax"
	case http.SameSiteNoneMode:
		return "None"
	default:
		return ""
	}
}

func getStringOrDefault(value interface{}, defaultValue string) string {
	if strValue, ok := value.(string); ok {
		return strValue
	}
	return defaultValue
}

func getBoolOrDefault(value interface{}, defaultValue bool) bool {
	if boolValue, ok := value.(bool); ok {
		return boolValue
	}
	return defaultValue
}
