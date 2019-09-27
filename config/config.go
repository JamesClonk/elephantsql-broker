package config

import (
	"strconv"
	"sync"

	"github.com/JamesClonk/elephantsql-broker/env"
)

var (
	config Config
	once   sync.Once
)

type Config struct {
	SkipSSL      string
	LogLevel     string
	LogTimestamp bool
	Username     string
	Password     string
	API          API
}
type API struct {
	URL                    string
	Key                    string
	DefaultRegion          string
	DefaultRegionPlansOnly bool
}

func loadConfig() {
	skipSSL, _ := strconv.ParseBool(env.Get("BROKER_SKIP_SSL_VALIDATION", "false"))
	logTimestamp, _ := strconv.ParseBool(env.Get("BROKER_LOG_TIMESTAMP", "false"))
	regionPlans, _ := strconv.ParseBool(env.Get("BROKER_API_DEFAULT_REGION_PLANS_ONLY", "true"))
	config = Config{
		SkipSSL:      skipSSL,
		LogLevel:     env.Get("BROKER_LOG_LEVEL", "info"),
		LogTimestamp: logTimestamp,
		Username:     env.MustGet("BROKER_USERNAME"),
		Password:     env.MustGet("BROKER_PASSWORD"),
		API: API{
			URL:                    env.Get("BROKER_API_URL", "https://customer.elephantsql.com/api"),
			Key:                    env.MustGet("BROKER_API_KEY"),
			DefaultRegion:          env.Get("BROKER_API_DEFAULT_REGION", "azure-arm::westeurope"),
			DefaultRegionPlansOnly: regionPlans,
		},
	}
}

func Get() *Config {
	configOnce.Do(func() {
		loadConfig()
	})
	return &config
}
