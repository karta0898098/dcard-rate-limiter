package serve

import (
	"context"
	"flag"
	"time"

	"github.com/karta0898098/dcard-rate-limiter/configs"
	restful "github.com/karta0898098/dcard-rate-limiter/pkg/delivery/http"
	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter"

	"github.com/karta0898098/kara/http"
	"github.com/karta0898098/kara/redis"
	"github.com/karta0898098/kara/zlog"

	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

// Run application
func Run() {
	var (
		path string
	)
	flag.StringVar(&path, "p", "", "serve -p ./deployments/config")
	flag.Parse()

	config := configs.NewConfig(path)

	app := fx.New(
		fx.Supply(config),
		ratelimiter.Module,
		fx.Provide(http.NewEcho),
		fx.Provide(redis.NewRedis),
		fx.Provide(restful.NewHandler),
		fx.Invoke(zlog.Setup),
		fx.Invoke(http.RunEcho),
		fx.Invoke(restful.SetupRoute),
	)
	app.Run()

	log.Info().Msg("Graceful Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		log.Info().Msgf("Server Shutdown: %v", err)
	}

	log.Info().Msg("Server exiting")
}
