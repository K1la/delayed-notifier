package config

import (
	"github.com/joho/godotenv"
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/zlog"
	"os"
)

const path = "./env"

func Init() *Config {
	wbCfg := config.New()

	err := wbCfg.Load(path)
	if err != nil {
		zlog.Logger.Panic().Err(err).Msg("could not read config file")
	}

	var cfg Config
	if err = wbCfg.Unmarshal(&cfg); err != nil {
		zlog.Logger.Panic().Err(err).Msg("could not unmarshal config file")
	}

	zlog.Logger.Info().Msgf("config: %+v", cfg)

	err = godotenv.Load(".env")
	if err != nil {
		zlog.Logger.Panic().Err(err).Msg("could not load .env file")
	}

	val, _ := os.LookupEnv("DB_PASSWORD")
	cfg.Postgres.Password = val

	return &cfg
}
