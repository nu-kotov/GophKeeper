package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_DefaultValues(t *testing.T) {
	cfg, err := NewConfig()
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, cfg.RunAddr, "localhost:8181", "unexpected RunAddr")
	assert.Equal(t, cfg.MiniIOConnection.Endpoint, "localhost:9000", "unexpected Endpoint")
	assert.Equal(t, cfg.SecretKey, "supersecretkey", "unexpected SecretKey")
	assert.Equal(t, cfg.TokenExp, time.Hour*72, "unexpected TokenExp")
	assert.Equal(t, cfg.MiniIOConnection.BucketName, "keeper", "unexpected MiniIOConnection.BucketName")
}
