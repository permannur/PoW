package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		Tcp      `yaml:"tcp"`
		Protocol `yaml:"protocol"`
		Pow      `yaml:"pow"`
		Log      `yaml:"logger"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"CLIENT_APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"CLIENT_APP_VERSION"`
	}

	Tcp struct {
		Port                 string `env-required:"true" yaml:"port" env:"CLIENT_TCP_PORT" env-default:"6677"`
		DialTimeoutSecond    int    `yaml:"dial_timeout_second" env:"CLIENT_DIAL_TIMEOUT_SECOND" env-default:"3"`
		ConnTimeoutSecond    int    `yaml:"conn_timeout_second" env:"CLIENT_CONN_TIMEOUT_SECOND" env-default:"5"`
		MaxDailCount         byte   `yaml:"max_dail_count_byte" env:"CLIENT_MAX_DAIL_COUNT_BYTE" env-default:"10"`
		DelayDailMillisecond int    `yaml:"delay_dail_millisecond" env:"CLIENT_DELAY_DAIL_MILLISECOND" env-default:"10"`
		ReadTimeoutSecond    int    `yaml:"read_timeout_second" env:"CLIENT_READ_TIMEOUT_SECOND" env-default:"3"`
	}

	Protocol struct {
		MaxReadSize uint32 `yaml:"max_read_size_uint32" env:"CLIENT_MAX_READ_SIZE_UINT32" env-default:"1024"`
	}

	Pow struct {
		MaxAttemptsCount uint32 `yaml:"max_attempts_count_uint32" env:"CLIENT_MAX_ATTEMPTS_COUNT_UINT32" env-default:"4294967295"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"CLIENT_LOG_LEVEL"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./client/config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config, read from file error: %s", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config, read from env error: %s", err)
	}

	return cfg, nil
}
