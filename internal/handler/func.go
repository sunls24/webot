package handler

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"log/slog"
	"strings"
	"webot/internal/cache"
	"webot/internal/constants"
	"webot/internal/types"
)

func isFunc(content string, cmd types.Func) bool {
	content = strings.TrimSpace(content)
	return len(content) == len(cmd) && strings.ToLower(content) == string(cmd)
}

func (h *Handler) onHelp(msg *openwechat.Message, ctx params) {
	_, err := msg.ReplyText(helpText(ctx.sender, ctx.main))
	if err != nil {
		slog.Error("reply help failed", slog.Any("err", err))
	}
}

func (h *Handler) onClear(msg *openwechat.Message, ctx params) {
	attr := slog.String("F", "onClear")

	err := h.DB.DeleteMessagesByAvatarID(ctx.main.AvatarID())
	if err != nil {
		slog.Error("delete messages failed", attr, slog.Any("err", err))
	}
	_, err = msg.ReplyText(ctx.atUser + clearText())
	if err != nil {
		slog.Error("reply clear failed", slog.Any("err", err))
	}
}

func (h *Handler) onToggle(msg *openwechat.Message, ctx params, key types.SettingsKey) {
	avatarID := ctx.main.AvatarID()
	value := cache.GetSettings(avatarID, key)
	if value == constants.ON {
		value = constants.OFF
	} else {
		value = constants.ON
	}
	cache.SetSettings(avatarID, key, value, h.DB)

	_, err := msg.ReplyText(ctx.atUser + v2exText(value))
	if err != nil {
		slog.Error(fmt.Sprintf("reply toggle %s failed", key), slog.Any("err", err))
	}
}
