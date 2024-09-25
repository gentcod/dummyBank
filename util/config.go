package util

import (
	"time"

	"github.com/gentcod/environ"
)

type Config struct {
	PortAddress string
	GrpcAddress string
	DBDriver string
	DBUrl string
	TokenSymmetricKey string
	TokenSecretKey string
	AccessTokenDuration time.Duration
	RefreshTokenDuration time.Duration
	MigrationUrl string
}

func LoadConfig(path string) (config Config, err error) {
	err = environ.Init(path, &config)
	if err != nil {
		return
	}

	return
}