package db

import (
	"os"
	"testing"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/hugocortes/hooks-api/bins/mocks"
	"github.com/hugocortes/hooks-api/bins/models"
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
var testCache *CacheRepo
var mockDB *mocks.DB
var mockBins []models.Bin
var accountID string

func TestCache(t *testing.T) {
	deps.LoadEnv()
	testRedisClient = deps.Redis()

	mockDB = new(mocks.DB)

	os.Setenv("REDIS_KEY", testPrefix)
	testCache = &CacheRepo{
		cache: deps.Redis(),
		db:    mockDB,
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

	t.Run("RedisShouldCacheIndividual", RedisShouldCacheIndividual)
}

func RedisShouldCacheIndividual(t *testing.T) {
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
