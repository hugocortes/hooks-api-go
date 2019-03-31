package cache

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultExpiration = "60000ms"
)

// Prefix ...
func Prefix() string {
	return os.Getenv("REDIS_KEY")
}

// Expiration ...
func Expiration() time.Duration {
	expire := os.Getenv("REDIS_EXPIRE_LOW") + "ms"
	if expire == "" {
		expire = defaultExpiration
	}

	duration, err := time.ParseDuration(expire)
	if err != nil {
		logrus.Panic("invalid cache expiration time: " + err.Error())
	}

	return duration
}

// GenKey ...
func GenKey(funcName string, args ...string) string {
	newString := Prefix() + ":"
	for _, arg := range args {
		newString += arg + ":"
	}
	return newString + funcName
}
