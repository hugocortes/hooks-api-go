package db

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hugocortes/hooks-api/bins/models"
	"github.com/hugocortes/hooks-api/common/deps"
	"github.com/hugocortes/hooks-api/migrations"
	gModels "github.com/hugocortes/hooks-api/models"
	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var db *gorm.DB
var testPostgres = &PostgresRepo{}

func TestPostgres(t *testing.T) {
	deps.LoadEnv()

	testDatabase := os.Getenv("POSTGRES_DB") + "-test"
	os.Setenv("POSTGRES_DB", testDatabase)

	db = deps.Postgres()
	testPostgres = &PostgresRepo{
		db: db,
	}

	migrations.Run(db)

	defer db.Close()

	t.Run("ShouldGetNoBins", ShouldGetNoBins)
	t.Run("ShouldCreateBin", ShouldCreateBin)
	t.Run("ShouldUpdateBin", ShouldUpdateBin)
	t.Run("ShouldGetBin", ShouldGetBin)
	t.Run("ShouldDeleteBin", ShouldDeleteBin)
	t.Run("ShouldDestroyBins", ShouldDestroyBins)
	t.Run("ShouldGetAllBins", ShouldGetAllBins)
}

func ShouldGetNoBins(t *testing.T) {
	accountID := uuid.New().String()
	opts := &gModels.QueryOpts{
		Page:  0,
		Limit: 10,
	}

	bins, err := testPostgres.GetAll(accountID, opts)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(bins))

	testPostgres.Destroy(accountID)
}

func ShouldCreateBin(t *testing.T) {
	accountID := uuid.New().String()
	title := fake.ProductName()
	bin := &models.Bin{
		Title:     title,
		AccountID: accountID,
	}

	binID, err := testPostgres.Create(bin)
	assert.Nil(t, err)
	assert.NotNil(t, binID, "BinID was not returned")

	testPostgres.Destroy(accountID)
}

func ShouldUpdateBin(t *testing.T) {
	accountID := uuid.New().String()
	createdBin := testCreateBin(accountID)

	createdBin.Title = "Updated Name"
	testPostgres.Update(accountID, createdBin.ID, createdBin)
	bin, err := testPostgres.Get(accountID, createdBin.ID)
	assert.Nil(t, err)
	assert.Equal(t, createdBin.Title, bin.Title)

	testPostgres.Destroy(accountID)
}

func ShouldGetBin(t *testing.T) {
	accountID := uuid.New().String()
	createdBin := testCreateBin(accountID)

	bin, err := testPostgres.Get(accountID, createdBin.ID)
	assert.Nil(t, err)
	assert.Equal(t, createdBin.ID, bin.ID)
	assert.Equal(t, accountID, bin.AccountID)
	assert.Equal(t, createdBin.Title, bin.Title)

	testPostgres.Destroy(accountID)
}

func ShouldDeleteBin(t *testing.T) {
	accountID := uuid.New().String()
	createdBin := testCreateBin(accountID)

	err := testPostgres.Delete(accountID, createdBin.ID)
	assert.Nil(t, err)

	bin, err := testPostgres.Get(accountID, createdBin.ID)
	assert.Nil(t, err)
	assert.Nil(t, bin)

	testPostgres.Destroy(accountID)
}

func ShouldDestroyBins(t *testing.T) {
	accountID := uuid.New().String()

	var bins []*models.Bin

	for i := 0; i < 10; i++ {
		bins = append(bins, testCreateBin(accountID))
	}

	testPostgres.Destroy(accountID)
	for i := 0; i < 10; i++ {
		bin, err := testPostgres.Get(accountID, bins[i].ID)
		assert.Nil(t, err)
		assert.Nil(t, bin)
	}
}

func ShouldGetAllBins(t *testing.T) {
	accountID := uuid.New().String()

	var bins []*models.Bin
	var err error

	for i := 0; i < 25; i++ {
		bins = append(bins, testCreateBin(accountID))
	}

	opts := &gModels.QueryOpts{Limit: -1, Page: -1}
	bins, err = testPostgres.GetAll(accountID, opts)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(bins))

	opts.Limit = 10
	opts.Page = 0
	bins, err = testPostgres.GetAll(accountID, opts)
	assert.Nil(t, err)
	assert.Equal(t, opts.Limit, len(bins))

	opts.Limit = 10
	opts.Page = 2
	bins, err = testPostgres.GetAll(accountID, opts)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(bins))

	testPostgres.Destroy(accountID)
}

func testCreateBin(accountID string) *models.Bin {
	title := fake.ProductName()
	bin := &models.Bin{
		Title:     title,
		AccountID: accountID,
	}
	binID, _ := testPostgres.Create(bin)
	bin.ID = binID
	return bin
}
