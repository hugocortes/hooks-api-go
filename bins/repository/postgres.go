package db

import (
	"errors"

	"github.com/hugocortes/hooks-api/bins/models"
	gModels "github.com/hugocortes/hooks-api/models"
	"github.com/jinzhu/gorm"
)

// PostgresRepo ...
type PostgresRepo struct {
	db *gorm.DB
}

// GetAll ...
func (r *PostgresRepo) GetAll(accountID string, opts gModels.QueryOpts) (*[]models.Bin, error) {
	return nil, errors.New("GetAll")
}

// Get ...
func (r *PostgresRepo) Get(accountID string, ID string) (*models.Bin, error) {
	return nil, errors.New("Get")
}

// Create ...
func (r *PostgresRepo) Create(bin *models.Bin) (string, error) {
	return "", errors.New("Create")
}

// Update ...
func (r *PostgresRepo) Update(accountID string, ID string, bin *models.Bin) error {
	return errors.New("Update")
}

// Delete ...
func (r *PostgresRepo) Delete(accountID string, ID string) error {
	return errors.New("Delete")
}

// Destroy ...
func (r *PostgresRepo) Destroy(accountID string) error {
	return errors.New("Destroy")
}
