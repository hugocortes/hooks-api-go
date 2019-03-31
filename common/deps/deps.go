package deps

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
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
// Returns true if it's in developer env, otherwise prod
func ConfigureLog() bool {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	switch os.Getenv("LOG_LEVEL") {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
		return true
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		return true
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		return false
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
		return false
	default:
		logrus.SetLevel(logrus.ErrorLevel)
		return false
	}
}

// Router creates the default gin router
func Router() *gin.Engine {
	gin.SetMode("release")
	if ConfigureLog() {
		gin.SetMode("debug")
	}

	return gin.Default()
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

	logrus.Info("redis connection √")

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

	db.SingularTable(true)

	logrus.Info("postgres connection √")

	return db
}
