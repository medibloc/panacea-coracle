package config_test

import (
	"testing"

	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMustLoad(t *testing.T) {
	t.Setenv("DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("DATAVAL_HTTP_LADDR", "0.0.0.0:8181")

	c := config.MustLoad()
	require.Equal(t, log.DebugLevel, log.Level(c.LogLevel))
	require.Equal(t, "0.0.0.0:8181", c.HTTPListenAddr)
}

func TestMustLoad_MissingRequiredEnv(t *testing.T) {
	t.Setenv("DATAVAL_LOG_LEVEL", "debug")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_InvalidLogLevel(t *testing.T) {
	t.Setenv("DATAVAL_LOG_LEVEL", "hello")
	t.Setenv("DATAVAL_HTTP_LADDR", "0.0.0.0:8181")

	require.Panics(t, func() {
		config.MustLoad()
	})
}
