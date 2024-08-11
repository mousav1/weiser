package storage

import (
	"fmt"

	"github.com/spf13/viper"
)

// InitializeStorage بارگذاری و ثبت Driverها
func InitializeStorage(defaultDriverName string, disksConfig *viper.Viper) error {
	// ثبت Driverها بر اساس پیکربندی
	for name, diskConfig := range disksConfig.AllSettings() {
		var driver Storage
		var err error

		disk := diskConfig.(map[string]interface{})
		switch disk["driver"] {
		case "local":
			driver = NewLocalDriver(disk["base_path"].(string))
		case "s3":
			driver, err = NewS3Driver(disk["bucket"].(string))
			if err != nil {
				return fmt.Errorf("failed to create S3 driver: %w", err)
			}
		default:
			return fmt.Errorf("unknown storage driver: %s", disk["driver"])
		}
		RegisterDriver(name, driver)
	}

	// انتخاب Driver پیش‌فرض
	if defaultDriverName == "" {
		return fmt.Errorf("default driver name is empty")
	}

	if err := SetDefaultDriver(defaultDriverName); err != nil {
		return fmt.Errorf("failed to select default storage driver: %w", err)
	}

	return nil
}
