package zookeeper

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/samuel/go-zookeeper/zk"
)

func NewZookeeper(config Config) (*zk.Conn, error) {
	c, _, err := zk.Connect([]string{config.Addr}, time.Second) // *10)
	if err != nil {
		log.Err(err).Msg("connect zookeeper occur error")
		return nil, err
	}

	log.Info().Msgf("connect to zookeeper %v success", config.Addr)

	return c, nil
}
