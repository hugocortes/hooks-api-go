package db_test

import (
	"os"
	"testing"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/hugocortes/hooks-api/bins/mocks"
	"github.com/hugocortes/hooks-api/bins/models"
	binsDB "github.com/hugocortes/hooks-api/bins/repository/db"
	"github.com/hugocortes/hooks-api/common/cache"
	"github.com/hugocortes/hooks-api/common/deps"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testPrefix = "hooks-api-test"
)

var testRedisClient *redis.Client
var testCache *binsDB.CacheRepo
var mockDB *mocks.DB
var mockBins []models.Bin
var accountID string

func testCacheSetup() {
	deps.LoadEnv("../../../.env")
	testRedisClient = deps.Redis()

	mockDB = new(mocks.DB)

	os.Setenv("REDIS_KEY", testPrefix)
	testCache = &binsDB.CacheRepo{
		Cache: deps.Redis(),
		DB:    mockDB,
	}

	accountID = uuid.New().String()
	mockBins = []models.Bin{
		models.Bin{
			ID:        uuid.New().String(),
			Title:     fake.Title(),
			AccountID: accountID,
		},
		models.Bin{
			ID:        uuid.New().String(),
			Title:     fake.Title(),
			AccountID: accountID,
		},
	}
}

func testCacheTearDown() {
	testCache.Cache.FlushAll()
}

func TestCachedGet(t *testing.T) {
	testCacheSetup()
	defer testCacheTearDown()

	bin := mockBins[0]

	var rawQueryCount = 0
	mockDB.On("Get", accountID, bin.ID).Return(&bin, nil).Run(func(args mock.Arguments) {
		rawQueryCount++
	})

	// raw query
	testCache.Get(bin.AccountID, bin.ID)
	// cached
	testCache.Get(bin.AccountID, bin.ID)
	testCache.Get(bin.AccountID, bin.ID)
	testCache.Get(bin.AccountID, bin.ID)

	testRedisClient.Del(cache.GenKey("Get", accountID, bin.ID))
	testCache.Get(bin.AccountID, bin.ID)
	testCache.Get(bin.AccountID, bin.ID)
	testCache.Get(bin.AccountID, bin.ID)

	assert.True(t, rawQueryCount == 2, "Query was called more than once")
}

func TestCacheInvalidation(t *testing.T) {
	testCacheSetup()
	defer testCacheTearDown()

	bin := mockBins[0]

	var rawQueryCount = 0
	var updated = false
	mockDB.On("Get", accountID, bin.ID).Return(&bin, nil).Run(func(args mock.Arguments) {
		rawQueryCount++
	})
	mockDB.On("Update", accountID, bin.ID, &bin).Return(1, nil).Run(func(args mock.Arguments) {
		updated = true
	})

	// raw query, rawQueryCount incremented
	testCache.Get(bin.AccountID, bin.ID)

	// cached response
	testCache.Get(bin.AccountID, bin.ID)
	testCache.Get(bin.AccountID, bin.ID)

	// cache invalidate
	bin.Title = fake.ProductName()
	testCache.Update(bin.AccountID, bin.ID, &bin)

	// raw query, rawQueryCount incremented
	testCache.Get(bin.AccountID, bin.ID)

	assert.True(t, updated, "Updated was not mocked")
	assert.True(t, rawQueryCount == 2, "Query was called more than once")
}
