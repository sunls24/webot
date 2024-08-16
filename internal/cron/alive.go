package cron

import (
	"fmt"
	"log/slog"
	"time"
	"webot/internal/bot"
	"webot/internal/context"
)

const (
	aliveSpec = "0 * * * *"
)

type alive context.Context

func (ctx alive) run() {
	attr := slog.String("F", "alive.run")

	if !bot.GetBot().Alive() {
		slog.Warn("bot is not alive", attr)
		return
	}

	self, err := bot.GetBot().GetCurrentUser()
	if err != nil {
		slog.Error("bot get current user failed", attr, slog.Any("err", err))
		return
	}
	_, err = self.FileHelper().SendText(fmt.Sprintf("I'm alive. [%d]", time.Now().Unix()))
	if err != nil {
		slog.Error("bot send alive failed", slog.Any("err", err))
	}
}
