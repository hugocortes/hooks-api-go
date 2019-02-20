package deps

import (
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// LoadEnv loads all env variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("did not load from .env")
	}
}

// ConfigureLog configures the logger
func ConfigureLog() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	switch os.Getenv("LOG_LEVEL") {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	}
}

// Redis returns a redis connection
func Redis() *redis.Client {
	options := &redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
	}

	redisAuth := os.Getenv("REDIS_AUTH")
	if redisAuth != "" {
		options.Password = redisAuth
	}

	redis := redis.NewClient(options)

	return redis
}

// Postgres creates the new postgres connection
func Postgres() *gorm.DB {
	pg := []string{
		"host=" + os.Getenv("POSTGRES_HOST"),
		"port=" + os.Getenv("POSTGRES_PORT"),
		"user=" + os.Getenv("POSTGRES_USER"),
		"dbname=" + os.Getenv("POSTGRES_DB"),
		"password=" + os.Getenv("POSTGRES_PASS"),
		"sslmode=" + os.Getenv("POSTGRES_SSL"),
	}

	db, err := gorm.Open("postgres", strings.Join(pg, " "))
	if err != nil {
		logrus.Fatal(err)
	}

	if level := os.Getenv("LOG_LEVEL"); level == "debug" || level == "trace" {
		db.LogMode(true)
	}

	logrus.Info("√ postgres connected")

	return db
}
