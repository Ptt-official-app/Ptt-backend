package config

import (
	"fmt"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/pelletier/go-toml"
)

type Config struct {
	BBSHome               string
	ListenPort            int16
	AccessTokenPrivateKey string
	AccessTokenPublicKey  string
	AccessTokenExpiresAt  time.Duration
}

var (
	logger = logging.NewLogger()
)

func NewDefaultConfig() (*Config, error) {
	return NewConfig("./conf/config_default.toml", "config.toml")
}

// NewConfig load and return global config from config files, it will load
// defaultPath first, then userPath. if defaultPath can not be read, it will
// return error. it userPath can not be read, it will ignore userPath.
// user configration will override default configuration.
func NewConfig(defaultPath, userPath string) (*Config, error) {

	var config *Config
	config = &Config{}
	logger.Debugf("load default config")

	defaultConfig, err := toml.LoadFile(defaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}
	applyConfig(config, defaultConfig)

	logger.Debugf("load user config")
	userConfig, err := toml.LoadFile(userPath)
	if err != nil {
		return config, nil
	}
	applyConfig(config, userConfig)

	return config, nil
}

func applyConfig(config *Config, rawConfig *toml.Tree) {
	logger.Debugf("apply rawConfig")
	var s string
	var i int64
	var ok bool

	s, ok = rawConfig.Get("bbs.home").(string)
	if ok {
		logger.Debugf("read rawConfig bbs.home: %s", s)
		config.BBSHome = s
	}
	i, ok = rawConfig.Get("networking.listen_port").(int64)
	if ok {
		logger.Debugf("read rawConfig networking.listen_port: %v", i)
		config.ListenPort = int16(i)
	}
	s, ok = rawConfig.Get("security.access_token_private_key").(string)
	if ok {
		logger.Debugf("read rawConfig security.access_token_private_key")
		config.AccessTokenPrivateKey = s
	}
	s, ok = rawConfig.Get("security.access_token_public_key").(string)
	if ok {
		logger.Debugf("read rawConfig security.access_token_public_key")
		config.AccessTokenPublicKey = s
	}

	s, ok = rawConfig.Get("security.access_token_expires_at").(string)
	if ok {
		var err error
		logger.Debugf("read rawConfig security.access_token_expires_at: %v", s)
		config.AccessTokenExpiresAt, err = time.ParseDuration(s)
		if err != nil {
			logger.Warningf("parse duration failed: %v", err)
		}
	}
}
