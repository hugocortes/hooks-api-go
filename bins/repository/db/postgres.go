package db

import (
	"github.com/google/uuid"
	"github.com/hugocortes/hooks-api/bins/models"
	gModels "github.com/hugocortes/hooks-api/models"
	"github.com/jinzhu/gorm"
)

const (
	tableName = "bin"
)

// PostgresRepo provides the database connection
type PostgresRepo struct {
	db *gorm.DB
}

// GetAll returns a page of bins for the account
func (r *PostgresRepo) GetAll(accountID string, opts *gModels.QueryOpts) ([]*models.Bin, error) {
	var bins []*models.Bin

	table := r.db.Table(tableName)
	table.Where("account_id = ?", accountID).Offset(opts.GetOffset()).Limit(opts.GetLimit()).Find(&bins)

	return bins, nil
}

// Get one bin associated with the given account id
func (r *PostgresRepo) Get(accountID string, ID string) (*models.Bin, error) {
	bin := &models.Bin{}

	table := r.db.Table(tableName)
	table.Where("id = ? AND account_id = ?", ID, accountID).Find(&bin)

	if !bin.Initialized() {
		bin = nil
	}

	return bin, nil
}

// Create inserts a new bin to the table and returns the ID
func (r *PostgresRepo) Create(bin *models.Bin) (string, error) {
	bin.ID = uuid.New().String()

	table := r.db.Table(tableName)
	table.Create(bin)

	return bin.ID, nil
}

// Update updates the bin with the provided values
func (r *PostgresRepo) Update(accountID string, ID string, bin *models.Bin) error {
	table := r.db.Table(tableName)
	table.Where("id = ? AND account_id = ?", ID, accountID).Update(&bin)

	return nil
}

// Delete removes a bin associated with the account
func (r *PostgresRepo) Delete(accountID string, ID string) error {
	table := r.db.Table(tableName)
	table.Where("id = ? AND account_id = ?", ID, accountID).Delete(&models.Bin{})

	return nil
}

// Destroy removes all bins associated with the account
func (r *PostgresRepo) Destroy(accountID string) error {
	table := r.db.Table(tableName)
	table.Where("account_id = ?", accountID).Delete(&models.Bin{})

	return nil
}
