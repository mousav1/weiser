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

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

// Connect initializes the database and Redis connections.
func Connect() (*gorm.DB, error) {
	if err := setupDatabases(); err != nil {
		return nil, err
	}
	if err := migrateAndSeed(); err != nil {
		return nil, err
	}
	return DB, nil
}

// setupDatabases configures and connects to the database and Redis.
func setupDatabases() error {
	var err error
	dbType := viper.GetString("database.db_type")

	switch dbType {
	case "mysql":
		DB, err = connectToMySQL()
	case "postgres":
		DB, err = connectToPostgres()
	default:
		return errors.New("unknown database type")
	}

	if err != nil {
		return errors.Wrap(err, "failed to connect to database")
	}

	// if err := setupRedis(); err != nil {
	// 	return err
	// }

	return nil
}

// connectToMySQL establishes a connection to the MySQL database.
func connectToMySQL() (*gorm.DB, error) {
	mysqlConfig := viper.GetStringMapString("database.mysql")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig["username"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"])

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to MySQL database")
	}

	if err := db.Exec("SET GLOBAL innodb_buffer_pool_size=1024 * 1024 * 1024").Error; err != nil {
		return nil, errors.Wrap(err, "failed to set temporary memory buffers for spatial indexes")
	}

	return db, nil
}

// connectToPostgres establishes a connection to the PostgreSQL database.
func connectToPostgres() (*gorm.DB, error) {
	postgresConfig := viper.GetStringMapString("database.postgres")
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

	if err := db.Exec("SET work_mem = '1GB'").Error; err != nil {
		return nil, errors.Wrap(err, "failed to set temporary memory buffers for spatial indexes")
	}

	return db, nil
}

// setupRedis configures and connects to Redis.
func setupRedis() error {
	redisConfig := viper.GetStringMapString("database.redis.default")
	var err error
	RedisClient, err = ConnectToRedis(redisConfig)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Redis")
	}
	return nil
}

// ConnectToRedis establishes a connection to Redis.
func ConnectToRedis(redisConfig map[string]string) (*redis.Client, error) {
	db, err := strconv.Atoi(redisConfig["dbname"])
	if err != nil {
		return nil, errors.Wrap(err, "invalid Redis DB number")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig["host"], redisConfig["port"]),
		Password: redisConfig["password"],
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, errors.Wrap(err, "failed to connect to Redis")
	}

	return client, nil
}

// WriteToRedis writes a key-value pair to Redis.
func WriteToRedis(key, value string) error {
	if RedisClient == nil {
		if err := setupRedis(); err != nil {
			return errors.Wrap(err, "failed to connect to Redis")
		}
	}

	if err := RedisClient.Set(context.Background(), key, value, 0).Err(); err != nil {
		return errors.Wrap(err, "failed to write value to Redis")
	}

	return nil
}

// ReadFromRedis reads a value from Redis by key.
func ReadFromRedis(key string) (string, error) {
	if RedisClient == nil {
		if err := setupRedis(); err != nil {
			return "", errors.Wrap(err, "failed to connect to Redis")
		}
	}

	value, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		return "", errors.Wrap(err, "failed to read value from Redis")
	}

	return value, nil
}

// Close closes database and Redis connections.
func Close() error {
	if err := closeDatabase(); err != nil {
		return err
	}
	if err := closeRedis(); err != nil {
		return err
	}
	return nil
}

// closeDatabase closes the database connection.
func closeDatabase() error {
	if DB != nil {
		db, err := DB.DB()
		if err != nil {
			return errors.Wrap(err, "failed to get database connection")
		}
		if err := db.Close(); err != nil {
			return errors.Wrap(err, "failed to close database connection")
		}
	}
	return nil
}

// closeRedis closes the Redis connection.
func closeRedis() error {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			return errors.Wrap(err, "failed to close Redis connection")
		}
	}
	return nil
}

// migrateAndSeed performs database migrations and seeds.
func migrateAndSeed() error {
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return errors.Wrap(err, "failed to migrate the schema")
	}
	return nil
}
