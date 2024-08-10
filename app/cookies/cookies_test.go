package cookies

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetCookie(t *testing.T) {
	viper.Set("cookie", map[string]interface{}{
		"name":     "mycookie",
		"path":     "/",
		"domain":   "",
		"expires":  "2030-12-31T00:00:00Z",
		"secure":   true,
		"samesite": 0,
		"httponly": true,
	})

	app := fiber.New()

	// Define the route
	app.Get("/", func(c *fiber.Ctx) error {
		SetCookie(c, "mycookie", "Hello, World!", time.Now().Add(time.Hour))
		return c.SendStatus(http.StatusOK)
	})

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Create a new Fiber context
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	// Validate response status
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Validate the cookie
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		t.Error("Expected at least one cookie, but found none")
	} else {
		assert.Equal(t, "mycookie", cookies[0].Name)
		assert.Equal(t, "Hello, World!", cookies[0].Value)
	}
}

func TestGetCookie(t *testing.T) {
	app := fiber.New()

	// Define the route
	app.Get("/", func(c *fiber.Ctx) error {
		value, err := GetCookie(c, "mycookie")
		if err != nil {
			return c.SendStatus(http.StatusNotFound)
		}
		return c.SendString(value)
	})

	// Create request and add cookie
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "mycookie",
		Value: "Hello, World!",
	})

	// Create a new Fiber context and handler
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Validate response status
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Validate response body
	assert.Equal(t, "Hello, World!", string(body))
}

func TestSetHttpCookie(t *testing.T) {
	name := "my_cookie"
	value := "cookie_value"
	expire := time.Hour
	secure := true
	httpOnly := true
	sameSite := http.SameSiteLaxMode

	// Create a new ResponseRecorder to capture the HTTP response
	w := httptest.NewRecorder()

	// Set the HTTP cookie
	SetHttpCookie(w, name, value, expire, secure, httpOnly, sameSite)

	// Get the cookies from the response
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == name {
			found = true
			if cookie.Value != value {
				t.Errorf("Invalid cookie value. Expected: %s, Got: %s", value, cookie.Value)
			}
			break
		}
	}

	if !found {
		t.Errorf("Cookie not found in ResponseWriter")
	}
}

func TestSameSiteToString(t *testing.T) {
	testCases := []struct {
		sameSite http.SameSite
		expected string
	}{
		{http.SameSiteDefaultMode, "Lax"},
		{http.SameSiteStrictMode, "Strict"},
		{http.SameSiteLaxMode, "Lax"},
		{http.SameSiteNoneMode, "None"},
		{http.SameSite(999), ""},
	}

	for _, tc := range testCases {
		result := sameSiteToString(tc.sameSite)
		if result != tc.expected {
			t.Errorf("Expected SameSiteToString(%v) to return '%s', but got '%s'", tc.sameSite, tc.expected, result)
		}
	}
}

func TestGetStringOrDefault(t *testing.T) {
	testCases := []struct {
		value          interface{}
		defaultValue   string
		expectedResult string
	}{
		{"test", "default", "test"},
		{42, "default", "default"},
		{nil, "default", "default"},
	}

	for _, tc := range testCases {
		result := getStringOrDefault(tc.value, tc.defaultValue)
		if result != tc.expectedResult {
			t.Errorf("Expected GetStringOrDefault(%v, %s) to return '%s', but got '%s'", tc.value, tc.defaultValue, tc.expectedResult, result)
		}
	}
}

func TestGetBoolOrDefault(t *testing.T) {
	testCases := []struct {
		value          interface{}
		defaultValue   bool
		expectedResult bool
	}{
		{true, false, true},
		{false, true, false},
		{"true", false, false},
		{nil, true, true},
	}

	for _, tc := range testCases {
		result := getBoolOrDefault(tc.value, tc.defaultValue)
		if result != tc.expectedResult {
			t.Errorf("Expected GetBoolOrDefault(%v, %t) to return '%t', but got '%t'", tc.value, tc.defaultValue, tc.expectedResult, result)
		}
	}
}
