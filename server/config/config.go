package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		Tcp      `yaml:"tcp"`
		Pow      `yaml:"pow"`
		Protocol `yaml:"protocol"`
		Log      `yaml:"logger"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"SERVER_APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"SERVER_APP_VERSION"`
	}

	Tcp struct {
		Port                   string `env-required:"true" yaml:"port" env:"SERVER_TCP_PORT"  env-default:"6677"`
		AcceptDelayMillisecond int    `yaml:"accept_delay_millisecond" env:"SERVER_ACCEPT_DELAY_MILLISECOND" env-default:"10"`
		ConnTimeoutSecond      int    `yaml:"conn_timeout_second" env:"SERVER_CONN_TIMEOUT" env-default:"5"`
		ReadTimeoutSecond      int    `yaml:"read_timeout_second" env:"SERVER_READ_TIMEOUT" env-default:"3"`
	}

	Pow struct {
		LeadingZeroCount byte `yaml:"leading_zero_count_byte" env:"SERVER_LEADING_ZERO_COUNT_BYTE" env-default:"20"`
	}

	Protocol struct {
		MaxReadSize uint32 `yaml:"max_read_size_uint32" env:"CLIENT_MAX_READ_SIZE_UINT32" env-default:"1024"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"SERVER_LOG_LEVEL"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./server/config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config, read from file error: %s", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config, read from env error: %s", err)
	}

	return cfg, nil
}
