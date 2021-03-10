package service

// Config define rate limiter config
type Config struct {
	MaxCount     int64 `mapstructure:"max_count"`
	RateLimitSec int64 `mapstructure:"rate_limit_sec"`
}
