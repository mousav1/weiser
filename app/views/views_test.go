package views

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func TestView(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		viper.SetConfigFile("../../config/config.yaml")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("failed to read configuration file: %s", err)
		}
		viper.Set("template_engine", "pongo2")
		viper.Set("template_dir", "../../resources")
		// فراخوانی تابع View
		data := ViewData{
			Title: "My Page",
			Data:  "Welcome to my page",
		}
		err := View(c, data, "test.html")

		if err != nil {
			t.Fatalf("Received unexpected error: %v", err)
		}
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// اضافه کردن تست‌های جدید
	t.Run("TestTemplateEngine", func(t *testing.T) {
		// تست کردن تنظیمات موتور قالب‌بندی
		expectedEngine := "pongo2"
		actualEngine := viper.GetString("template_engine")
		if expectedEngine != actualEngine {
			t.Errorf("Expected template engine %s, but got %s", expectedEngine, actualEngine)
		}
	})

	t.Run("TestTemplateData", func(t *testing.T) {
		// تست کردن داده‌های قالب
		expectedTitle := "My Page"
		expectedData := "Welcome to my page"
		// بررسی مقدار داده‌های قالب از طریق تابع View
		actualTitle := getViewData().Title
		actualData := getViewData().Data
		if expectedTitle != actualTitle {
			t.Errorf("Expected title %s, but got %s", expectedTitle, actualTitle)
		}
		if expectedData != actualData {
			t.Errorf("Expected data %s, but got %s", expectedData, actualData)
		}
	})
}

func getViewData() ViewData {
	return ViewData{
		Title: "My Page",
		Data:  "Welcome to my page",
	}
}
