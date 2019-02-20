package db

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/hugocortes/hooks-api/bins"
	"github.com/hugocortes/hooks-api/bins/models"
	gModels "github.com/hugocortes/hooks-api/models"
)

// CacheRepo ..
type CacheRepo struct {
	db    bins.DB
	cache *redis.Client
}

// GetAll ...
func (r *CacheRepo) GetAll(accountID string, opts gModels.QueryOpts) (*[]models.Bin, error) {
	return nil, errors.New("GetAll")
}

// Get ...
func (r *CacheRepo) Get(accountID string, ID string) (*models.Bin, error) {
	return nil, errors.New("Get")
}

// Create ...
func (r *CacheRepo) Create(bin *models.Bin) (string, error) {
	return "", errors.New("Create")
}

// Update ...
func (r *CacheRepo) Update(accountID string, ID string, bin *models.Bin) error {
	return errors.New("Update")
}

// Delete ...
func (r *CacheRepo) Delete(accountID string, ID string) error {
	return errors.New("Delete")
}

// Destroy ...
func (r *CacheRepo) Destroy(accountID string) error {
	return errors.New("Destroy")
}
