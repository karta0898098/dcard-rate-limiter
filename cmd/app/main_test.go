package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/karta0898098/dcard-rate-limiter/configs"
	restful "github.com/karta0898098/dcard-rate-limiter/pkg/delivery/http"
	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter"
	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/service"
	httpInfra "github.com/karta0898098/kara/http"

	"github.com/karta0898098/dcard-rate-limiter/internal/zookeeper"
	rc "github.com/karta0898098/kara/redis"
	"github.com/karta0898098/kara/zlog"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	"go.uber.org/fx"
)

type rateLimiterTestSuite struct {
	suite.Suite
	app        *fx.App
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
		Out: fx.Out{},
		Log: zlog.Config{
			Env:   "test",
			AppID: "rate_limiter",
			Debug: true,
		},
		HTTP: httpInfra.Config{
			Mode: "debug",
			Port: ":18080",
		},
		Redis: rc.Config{
			Address: "127.0.0.1:16379",
		},
		RateLimit: service.Config{
			MaxCount:     60,
			RateLimitSec: 60,
		},
		Zookeeper: zookeeper.Config{
			Addr: "127.0.0.1:12181",
		},
	}

	s.app = fx.New(
		fx.Supply(config),
		ratelimiter.Module,
		fx.Provide(zookeeper.NewZookeeper),
		fx.Provide(rc.NewRedis),
		fx.Provide(httpInfra.NewEcho),
		fx.Provide(restful.NewHandler),
		fx.Invoke(zlog.Setup),
		fx.Invoke(httpInfra.RunEcho),
		fx.Invoke(restful.SetupRoute),
		fx.Populate(&s.redisClint),
	)
	go s.app.Run()
}

func (s *rateLimiterTestSuite) SetupTest() {
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

func (s *rateLimiterTestSuite) TestRateLimiterGotError() {
	var (
		task    int
		wg      sync.WaitGroup
		success int
		failed  int
	)

	task = 61
	wg.Add(task)

	for i := 0; i < task; i++ {
		go func(i int) {
			defer wg.Done()

			client := &http.Client{}
			req, _ := http.NewRequest("GET", "http://localhost:18080/api/v1/protected", nil)
			req.Header.Add(echo.HeaderXRealIP, "127.0.0.1")
			resp, _ := client.Do(req)
			data, _ := ioutil.ReadAll(resp.Body)
			log.Debug().Msg(string(data))
			if resp.StatusCode == 429 {
				failed++
			}

			if resp.StatusCode == 200 {
				success++
			}
		}(i)
	}
	wg.Wait()

	s.Equal(60, success)
	s.Equal(1, failed)
}
