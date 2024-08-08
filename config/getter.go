package config

import (
	"github.com/eatmoreapple/openwechat"
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
