package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	errors "github.com/mousav1/weiser/app/error"
	"github.com/mousav1/weiser/app/models"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client

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

func ConnectToRedis(redisConfig map[string]string) (*redis.Client, error) {
	// Load Redis settings from configuration file
	db, err := strconv.Atoi(redisConfig["dbname"])
	if err != nil {
		return nil, err
	}
	// Connect to Redis using the settings loaded from the configuration file
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig["host"] + ":" + redisConfig["port"],
		Password: redisConfig["password"],
		DB:       db, // use default DB
	})

	// Test the connection to Redis
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to Redis")
	}

	return client, nil
}

func WriteToRedis(key string, value string) error {
	// Connect to Redis
	if RedisClient == nil {
		var err error
		redisConfig := viper.GetStringMapString("database.redis.default")
		RedisClient, err = ConnectToRedis(redisConfig)
		if err != nil {
			return errors.Wrap(err, "failed to connect to Redis")
		}
	}

	// Write the value to Redis
	err := RedisClient.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return errors.Wrap(err, "failed to write value to Redis")
	}

	return nil
}

func ReadFromRedis(key string) (string, error) {
	// Connect to Redis
	if RedisClient == nil {
		var err error
		redisConfig := viper.GetStringMapString("database.redis.default")
		RedisClient, err = ConnectToRedis(redisConfig)
		if err != nil {
			return "", errors.Wrap(err, "failed to connect to Redis")
		}
	}

	// Read the value from Redis
	value, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		return "", errors.Wrap(err, "failed to read value from Redis")
	}

	return value, nil
}

func Close() error {
	// Close the database connection
	if DB != nil {
		db, err := DB.DB()
		if err != nil {
			return errors.Wrap(err, "failed to get database connection")
		}
		err = db.Close()
		if err != nil {
			return errors.Wrap(err, "failed to close database connection")
		}
	}

	// Close the Redis connection
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			return errors.Wrap(err, "failed to close Redis connection")
		}
	}

	return nil
}

func migrateAndSeed() error {
	// Migrate the schema
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		return errors.Wrap(err, "failed to migrate the schema")
	}

	return nil
}
