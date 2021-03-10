package main

import (
	"context"
	"flag"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/karta0898098/dcard-rate-limiter/configs"
	restful "github.com/karta0898098/dcard-rate-limiter/pkg/delivery/http"
	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter"
	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/service"

	rc "github.com/karta0898098/kara/redis"
	"github.com/karta0898098/kara/tracer"
	"github.com/karta0898098/kara/zlog"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	"go.uber.org/fx"
)

type rateLimiterTestSuite struct {
	suite.Suite
	app        *fx.App
	handler    *restful.Handler
	c          echo.Context
	resp       *httptest.ResponseRecorder
	redisClint *redis.Client
}

func TestRateLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(rateLimiterTestSuite))
}

func (s *rateLimiterTestSuite) SetupSuite() {
	var (
		path string
	)
	flag.StringVar(&path, "p", "", "serve -p ./deployments/config")
	flag.Parse()

	config := configs.Configurations{
		Log: zlog.Config{
			Env:   "test",
			AppID: "rate_limiter",
			Debug: true,
		},
		Redis: rc.Config{
			Address: "127.0.0.1:16379",
		},
		RateLimit: service.Config{
			MaxCount:     60,
			RateLimitSec: 60,
		},
	}

	s.app = fx.New(
		fx.Supply(config),
		ratelimiter.Module,
		fx.Provide(rc.NewRedis),
		fx.Provide(restful.NewHandler),
		fx.Invoke(zlog.Setup),
		fx.Populate(&s.handler),
		fx.Populate(&s.redisClint),
	)
	go s.app.Run()
}

func (s *rateLimiterTestSuite) SetupTest() {
	e := echo.New()

	req := httptest.NewRequest("GET", "/proctored", nil)
	req.Header.Set(echo.HeaderXRealIP, "127.0.0.1")
	s.resp = httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), tracer.TraceIDKey, uuid.New().String())
	s.c = e.NewContext(req, s.resp)
	s.c.SetPath("/proctored")
	s.c.SetRequest(s.c.Request().WithContext(ctx))
}

func (s *rateLimiterTestSuite) TearDownTest() {
	log.Info().Msg("Graceful Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	s.redisClint.FlushAll(ctx)
	defer cancel()
	if err := s.app.Stop(ctx); err != nil {
		log.Info().Msgf("Server Shutdown: %v", err)
	}

	log.Info().Msg("Server exiting")
}

func (s *rateLimiterTestSuite) TestRateLimiter() {
	_ = s.handler.ProtectedEndpoint(s.c)
	s.Equal(s.resp.Code, 200)
}

func (s *rateLimiterTestSuite) TestRateLimiterGotError() {
	var (
		task int
	)

	task = 60
	for i := 0; i < task; i++ {
		err := s.handler.ProtectedEndpoint(s.c)
		if i > task {
			s.NotNil(err)
		} else {
			s.NoError(err)
		}
	}
}
