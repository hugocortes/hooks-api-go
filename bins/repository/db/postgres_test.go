package db_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hugocortes/hooks-api/bins/models"
	binsDB "github.com/hugocortes/hooks-api/bins/repository/db"
	"github.com/hugocortes/hooks-api/common/deps"
	"github.com/hugocortes/hooks-api/migrations"
	gModels "github.com/hugocortes/hooks-api/models"
	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var db *gorm.DB
var testPostgres = &binsDB.PostgresRepo{}
var originalDb string

func testPostgresSetup() {
	deps.LoadEnv("../../../.env")
	originalDb = os.Getenv("POSTGRES_DB")
	testDatabase := originalDb + "-test"
	os.Setenv("POSTGRES_DB", testDatabase)

	db = deps.Postgres()
	testPostgres = &binsDB.PostgresRepo{
		DB: db,
	}

	migrations.Run(db)
}

func testPostgresTearDown() {
	os.Setenv("POSTGRES_DB", originalDb)
	db.Close()
}

func TestGetNoBins(t *testing.T) {
	testPostgresSetup()
	defer testPostgresTearDown()

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

func TestCreateBin(t *testing.T) {
	testPostgresSetup()
	defer testPostgresTearDown()

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

func TestUpdateBin(t *testing.T) {
	testPostgresSetup()
	defer testPostgresTearDown()

	accountID := uuid.New().String()
	createdBin := testCreateBin(accountID)
	updatedTs := createdBin.UpdatedAt

	time.Sleep(2 * time.Second)

	createdBin.Title = "Updated Name"
	testPostgres.Update(accountID, createdBin.ID, createdBin)
	bin, err := testPostgres.Get(accountID, createdBin.ID)
	assert.Nil(t, err)
	assert.Equal(t, createdBin.Title, bin.Title)
	assert.Equal(t, createdBin.CreatedAt.Unix(), bin.CreatedAt.Unix())
	assert.NotEqual(t, updatedTs.UTC(), bin.UpdatedAt.UTC())

	testPostgres.Destroy(accountID)
}

func TestGetBin(t *testing.T) {
	testPostgresSetup()
	defer testPostgresTearDown()

	accountID := uuid.New().String()
	createdBin := testCreateBin(accountID)

	bin, err := testPostgres.Get(accountID, createdBin.ID)
	assert.Nil(t, err)
	assert.Equal(t, createdBin.ID, bin.ID)
	assert.Equal(t, accountID, bin.AccountID)
	assert.Equal(t, createdBin.Title, bin.Title)

	testPostgres.Destroy(accountID)
}

func TestDeleteBin(t *testing.T) {
	testPostgresSetup()
	defer testPostgresTearDown()

	accountID := uuid.New().String()
	createdBin := testCreateBin(accountID)

	err := testPostgres.Delete(accountID, createdBin.ID)
	assert.Nil(t, err)

	bin, err := testPostgres.Get(accountID, createdBin.ID)
	assert.NotNil(t, err, "Expected a not found error")
	assert.Nil(t, bin)

	testPostgres.Destroy(accountID)
}

func TestDestroyBins(t *testing.T) {
	testPostgresSetup()
	defer testPostgresTearDown()

	accountID := uuid.New().String()

	var bins []*models.Bin

	for i := 0; i < 10; i++ {
		bins = append(bins, testCreateBin(accountID))
	}

	testPostgres.Destroy(accountID)
	for i := 0; i < 10; i++ {
		bin, err := testPostgres.Get(accountID, bins[i].ID)
		assert.NotNil(t, err, "Expected not found error")
		assert.Nil(t, bin)
	}
}

func TestGetAllBins(t *testing.T) {
	testPostgresSetup()
	defer testPostgresTearDown()

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

func TestGetErrors(t *testing.T) {
	testPostgresSetup()
	defer testPostgresTearDown()

	binID := uuid.New().String()
	accountID := uuid.New().String()

	var bin *models.Bin
	var err error

	bin, err = testPostgres.Get(accountID, binID)
	assert.NotNil(t, err, "Expected an error")
	assert.Nil(t, bin, "Expected nil")

	err = testPostgres.Update(accountID, binID, &models.Bin{})
	assert.NotNil(t, err, "Expected an error")

	err = testPostgres.Delete(accountID, binID)
	assert.NotNil(t, err, "Expected an error")

	err = testPostgres.Destroy(accountID)
	assert.NotNil(t, err, "Expected an error")
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
