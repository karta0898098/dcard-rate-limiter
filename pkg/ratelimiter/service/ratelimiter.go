package service

import (
	"context"

	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/domain"
	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/repository"

	"github.com/karta0898098/kara/errors"
)

var _ domain.RateLimiterService = &rateLimiterService{}

type rateLimiterService struct {
	maxCount       int64
	timeUintSecond int64
	repo           repository.RateLimiterRepository
}

// NewRateLimiterService ...
func NewRateLimiterService(repo repository.RateLimiterRepository, config Config) domain.RateLimiterService {
	return &rateLimiterService{
		maxCount:       config.MaxCount,
		timeUintSecond: config.RateLimitSec,
		repo:           repo,
	}
}

// RequireResource require access resource
// if too many request will return error
func (srv *rateLimiterService) RequireResource(ctx context.Context, addr string, url string) (claims *domain.Claims, err error) {
	var (
		key string
	)

	if addr == "" {
		return nil, errors.ErrInternal.Build("addr must need input")
	}

	// why key format is url + : + addr ?
	// because I want to distinguish url require count
	// every url resource can has own rate limit
	key = url + ":" + addr

	count, err := srv.repo.GetRequestCount(ctx, key)
	if err != nil {
		return nil, err
	}

	if count >= srv.maxCount {
		return nil, errors.ErrTooManyRequests.Build("too many request please retry later addr = %s url = %s", addr, url)
	}

	err = srv.repo.AddRequestCount(ctx, key, srv.timeUintSecond)
	if err != nil {
		return nil, err
	}

	// why count need add 1
	// because first input count is zero
	return &domain.Claims{
		URL:          url,
		RemoteAddr:   addr,
		RequestCount: count + 1,
	}, nil
}
