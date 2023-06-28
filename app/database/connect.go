package database

import (
	"fmt"

	"github.com/mousav1/weiser/app/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func connectToDB() (*gorm.DB, error) {
//     dbConfig := viper.GetStringMapString("database")
//     dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig["username"], dbConfig["password"], dbConfig["host"], dbConfig["port"], dbConfig["dbname"])
//     db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//     if err != nil {
//         return nil, err
//     }

//     return db, nil
// }

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	// Load database settings from configuration file
	dbType := viper.GetString("database.db_type")

	var err error

	// Connect to the database based on the type specified in the configuration file
	switch dbType {
	case "mysql":
		DB, err = connectToMySQL()
	case "postgres":
		DB, err = connectToPostgres()
	default:
		return nil, fmt.Errorf("Unknown database type")
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %s", err)
	}

	// Set up database migrations and seeds
	if err := migrateAndSeed(); err != nil {
		return nil, fmt.Errorf("Failed to run migrations and seeds: %s", err)
	}

	return DB, nil
}

// Connect to MySQL database
func connectToMySQL() (*gorm.DB, error) {

	// Load database settings from configuration file
	mysqlConfig := viper.GetStringMapString("database.mysql")

	// Connect to the database using the settings loaded from the configuration file
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig["username"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"])), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Set temporary memory buffers for spatial indexes
	db.Exec("SET GLOBAL innodb_buffer_pool_size=1024 * 1024 * 1024")

	return db, nil
}

// Connect to PostgreSQL database
func connectToPostgres() (*gorm.DB, error) {
	// Load database settings from configuration file
	postgresConfig := viper.GetStringMapString("database.postgres")

	// Connect to the database using the settings loaded from the configuration file
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		postgresConfig["host"],
		postgresConfig["user"],
		postgresConfig["password"],
		postgresConfig["dbname"],
		postgresConfig["port"])), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Set temporary memory buffers for spatial indexes
	db.Exec("SET work_mem = '1GB'")

	return db, nil
}

func migrateAndSeed() error {
	// Migrate the schema
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	// Create a user if there are no users in the database
	// var count int64
	// if err := DB.Model(&models.User{}).Count(&count).Error; err != nil {
	// 	return err
	// }

	// if count == 0 {
	// 	user := &models.User{Name: "John Doe", Email: "john.doe@example.com", Age: 30}
	// 	if err := DB.Create(user).Error; err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
