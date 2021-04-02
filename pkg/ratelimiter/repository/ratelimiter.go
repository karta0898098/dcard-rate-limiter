package repository

import (
	"context"

	"github.com/karta0898098/kara/errors"
	"github.com/karta0898098/kara/timeutil"
	"github.com/karta0898098/kara/tracer"

	"github.com/go-redis/redis/v8"
)

var _ RateLimiterRepository = &rateLimiterRepository{}

// RateLimiterRepository define rate limit  method
type RateLimiterRepository interface {
	// AddRequestCount implement redis slide windows algorithm
	// When too many request will return -1
	// ref : https://en.wikipedia.org/wiki/Sliding_window_protocol
	AddRequestCount(ctx context.Context, key string, maxCount, timeUintSecond int64) (count int64, err error)
}

type rateLimiterRepository struct {
	redisClient *redis.Client
	script      string
}

// NewRateLimiterRepository ...
func NewRateLimiterRepository(redisClient *redis.Client) RateLimiterRepository {
	script := `
	local key = KEYS[1]
	local maxCount = tonumber(ARGV[1])
	local nowTimeMS = tonumber(ARGV[2])
	local maxScoreMS = tonumber(ARGV[3])
	local timeUintSecond = tonumber(ARGV[4])
	local traceID = ARGV[5]

	local count = redis.call('ZCARD', key)
	if count >= maxCount then
    	return -1
	end

	redis.call('ZREMRANGEBYSCORE', key, 0, maxScoreMS)
	redis.call('ZADD', key, nowTimeMS, traceID)
	redis.call('EXPIRE', key, timeUintSecond)

	return count`
	return &rateLimiterRepository{
		redisClient: redisClient,
		script:      script,
	}
}

// AddRequestCount implement redis slide windows algorithm
// When too many request will return -1
// ref : https://en.wikipedia.org/wiki/Sliding_window_protocol
func (repo *rateLimiterRepository) AddRequestCount(ctx context.Context, key string, maxCount, timeUintSecond int64) (count int64, err error) {
	nowTimeMS := timeutil.NowMS()
	maxScoreMS := nowTimeMS - (timeUintSecond * 1000)
	traceID := ctx.Value(tracer.TraceIDKey)

	args := []interface{}{maxCount, nowTimeMS, maxScoreMS, timeUintSecond, traceID}
	script := redis.NewScript(repo.script)
	count, err = script.Run(ctx, repo.redisClient, []string{key}, args...).Int64()
	if err != nil {
		return 0, errors.ErrInternal.Build("exec redis occur error reason :%v", err)
	}

	return count, nil
}
