package main

import (
	"github.com/pelletier/go-toml"
	"time"
)

type GlobalConfig struct {
	BBSHome               string
	ListenPort            int16
	AccessTokenPrivateKey string
	AccessTokenExpiresAt  time.Duration
}

var globalConfig = GlobalConfig{}

func loadDefaultConfig() {
	defaultConfig, err := toml.LoadFile("config_default.toml")
	if err != nil {
		logger.Errorf("load default config error: %v", err)
		return
	}
	applyConfig(defaultConfig)
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		logger.Warningf("load custom config error: %v", err)
		return
	}
	applyConfig(config)

}

func applyConfig(config *toml.Tree) {
	var s string
	var i int64
	var ok bool

	s, ok = config.Get("bbs.home").(string)
	if ok {
		globalConfig.BBSHome = s
	}
	i, ok = config.Get("networking.listen_port").(int64)
	if ok {
		globalConfig.ListenPort = int16(i)
	}
	s, ok = config.Get("security.access_token_private_key").(string)
	if ok {
		globalConfig.AccessTokenPrivateKey = s
	}
	s, ok = config.Get("security.access_token_expires_at").(string)
	if ok {
		var err error
		globalConfig.AccessTokenExpiresAt, err = time.ParseDuration(s)
		if err != nil {
			logger.Warningf("parse duartion failed: %v", err)
		}
	}

}
