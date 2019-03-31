package repository

import (
	"github.com/go-redis/redis"
	"github.com/hugocortes/hooks-api/bins"
	"github.com/hugocortes/hooks-api/bins/repository/db"
	"github.com/jinzhu/gorm"
)

// New ...
func New(postgres *gorm.DB, redis *redis.Client) *bins.Repository {
	return &bins.Repository{
		DB: db.New(postgres, redis),
	}
}
