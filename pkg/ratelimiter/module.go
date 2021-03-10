package ratelimiter

import (
	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/repository"
	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/service"

	"go.uber.org/fx"
)

// Module provide uber fx
var Module = fx.Options(
	fx.Provide(service.NewRateLimiterService),
	fx.Provide(repository.NewRateLimiterRepository),
)
