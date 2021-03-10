package domain

import "context"

// Claims define rate limiter model
type Claims struct {
	URL          string `json:"url"`
	RemoteAddr   string `json:"remoteAddr"`
	RequestCount int64  `json:"requestCount"`
}

// RateLimiterService define rate limiter
type RateLimiterService interface {
	// RequireResource try to get claims
	RequireResource(ctx context.Context, addr string, url string) (claims *Claims, err error)
}
