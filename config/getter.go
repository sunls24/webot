package config

import (
	"github.com/eatmoreapple/openwechat"
	"io"
	"webot/internal/types"
)

func (cfg *Config) GetMode() openwechat.BotPreparer {
	switch cfg.Mode {
	case types.LoginDesktop:
		return openwechat.Desktop
	default:
		return openwechat.Normal
	}
}

func (cfg *Config) NewStorage() io.ReadWriteCloser {
	return openwechat.NewFileHotReloadStorage(cfg.Storage)
}

func (cfg *Config) GetModel(slow bool) string {
	if slow {
		return cfg.OpenAI.SlowModel
	}
	return cfg.OpenAI.Model
}
