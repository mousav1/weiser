package database

import (
	"context"
	"fmt"

	errors "github.com/mousav1/weiser/app/error"
	"github.com/mousav1/weiser/app/models"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
		return nil, errors.New("unknown database type")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	// Set up database migrations and seeds
	if err := migrateAndSeed(); err != nil {
		return nil, errors.Wrap(err, "failed to run migrations and seeds")
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
		return nil, errors.Wrap(err, "failed to connect to MySQL database")
	}

	// Set temporary memory buffers for spatial indexes
	if err := db.Exec("SET GLOBAL innodb_buffer_pool_size=1024 * 1024 * 1024").Error; err != nil {
		return nil, errors.Wrap(err, "failed to set temporary memory buffers for spatial indexes")
	}

	return db, nil
}

// Connect to PostgreSQL database
func connectToPostgres() (*gorm.DB, error) {
	// Load database settings from configuration file
	postgresConfig := viper.GetStringMapString("database.postgres")

	// Connect to the database using the settings loaded from the configuration file
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		postgresConfig["host"],
		postgresConfig["user"],
		postgresConfig["password"],
		postgresConfig["dbname"],
		postgresConfig["port"],
		postgresConfig["sslmode"])
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to PostgreSQL database")
	}

	// Set temporary memory buffers for spatial indexes
	if err := db.Exec("SET work_mem = '1GB'").Error; err != nil {
		return nil, errors.Wrap(err, "failed to set temporary memory buffers for spatial indexes")
	}

	return db, nil
}

// Connect to Redis database
func ConnectToRedis() (*redis.Client, error) {
	// Load Redis settings from configuration file
	redisConfig := viper.GetStringMapString("database.redis")

	// Connect to the Redis database using the settings loaded from the configuration file
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig["host"], redisConfig["port"]),
		Password: redisConfig["password"],
		DB:       0,
	})

	// Test the connection to the Redis database
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to Redis database")
	}

	return client, nil
}

func migrateAndSeed() error {
	// Migrate the schema
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		return errors.Wrap(err, "failed to migrate the schema")
	}

	return nil
}
