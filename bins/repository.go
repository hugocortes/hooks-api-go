package bins

import (
	"github.com/hugocortes/hooks-api/bins/models"
	gModels "github.com/hugocortes/hooks-api/models"
)

// Repository ...
type Repository struct {
	DB DB
}

// DB ...
type DB interface {
	GetAll(accountID string, opts *gModels.QueryOpts) ([]*models.Bin, error)
	Get(accountID string, ID string) (*models.Bin, error)
	Create(bin *models.Bin) (string, error)
	Update(accountID string, ID string, bin *models.Bin) error
	Delete(accountID string, ID string) error
	Destroy(accountID string) error
}
