package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/karta0898098/kara/errors"
	"github.com/karta0898098/kara/timeutil"
	"github.com/karta0898098/kara/tracer"
)

var _ RateLimiterRepository = &rateLimiterRepository{}

type RateLimiterRepository interface {
	AddRequestCount(ctx context.Context, key string, timeUintSecond int64) (err error)
	GetRequestCount(ctx context.Context, key string) (count int64, err error)
}

type rateLimiterRepository struct {
	redisClient *redis.Client
}

func NewRateLimiterRepository(redisClient *redis.Client) RateLimiterRepository {
	return &rateLimiterRepository{
		redisClient: redisClient,
	}
}

func (repo *rateLimiterRepository) AddRequestCount(ctx context.Context, key string, timeUintSecond int64) (err error) {
	var (
		nowTimeMS  int64
		maxScoreMS int64
		traceID    string
	)

	nowTimeMS = timeutil.NowMS()
	maxScoreMS = nowTimeMS - (timeUintSecond * 1000)
	traceID = ctx.Value(tracer.TraceIDKey).(string)

	pipe := repo.redisClient.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(maxScoreMS, 10))
	pipe.ZAdd(ctx, key, &redis.Z{
		Score:  float64(nowTimeMS),
		Member: traceID,
	})
	pipe.Expire(ctx, key, time.Duration(timeUintSecond)*time.Second)
	_, err = pipe.Exec(ctx)

	if err != nil {
		return errors.ErrInternal.BuildWithError(err)
	}

	return nil
}

func (repo *rateLimiterRepository) GetRequestCount(ctx context.Context, key string) (count int64, err error) {
	count, err = repo.redisClient.ZCard(ctx, key).Result()
	if err != nil {
		return 0, errors.ErrInternal.BuildWithError(err)
	}
	return count, nil
}
