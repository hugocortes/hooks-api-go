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
	DB *gorm.DB
}

// GetAll returns a page of bins for the account
func (r *PostgresRepo) GetAll(accountID string, opts *gModels.QueryOpts) ([]*models.Bin, error) {
	var bins []*models.Bin

	table := r.DB.Table(tableName)
	table.Where("account_id = ?", accountID).Offset(opts.GetOffset()).Limit(opts.GetLimit()).Find(&bins)

	return bins, nil
}

// Get one bin associated with the given account id
func (r *PostgresRepo) Get(accountID string, ID string) (*models.Bin, error) {
	bin := &models.Bin{}

	table := r.DB.Table(tableName)
	res := table.Where("id = ? AND account_id = ?", ID, accountID).Find(&bin).RecordNotFound()

	if res {
		bin = nil
	}

	return bin, nil
}

// Create inserts a new bin to the table
func (r *PostgresRepo) Create(accountID string, bin *models.Bin) error {
	bin.AccountID = accountID
	bin.ID = uuid.New().String()

	table := r.DB.Table(tableName)
	table.Create(&bin)

	return nil
}

// Update updates the bin with the provided values
func (r *PostgresRepo) Update(accountID string, ID string, bin *models.Bin) (int, error) {
	table := r.DB.Table(tableName)

	res := table.Model(&models.Bin{}).Where("id = ? AND account_id = ?", ID, accountID).Omit("created_at", "id", "account_id").Update(&bin)

	return int(res.RowsAffected), nil
}

// Delete removes a bin associated with the account
func (r *PostgresRepo) Delete(accountID string, ID string) (int, error) {
	table := r.DB.Table(tableName)
	res := table.Where("id = ? AND account_id = ?", ID, accountID).Delete(&models.Bin{})

	return int(res.RowsAffected), nil
}

// Destroy removes all bins associated with the account
func (r *PostgresRepo) Destroy(accountID string) (int, error) {
	table := r.DB.Table(tableName)
	res := table.Where("account_id = ?", accountID).Delete(&models.Bin{})

	return int(res.RowsAffected), nil
}
