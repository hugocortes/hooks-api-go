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
	Create(accountID string, bin *models.Bin) error
	Update(accountID string, ID string, bin *models.Bin) (int, error)
	Delete(accountID string, ID string) (int, error)
	Destroy(accountID string) (int, error)
}
