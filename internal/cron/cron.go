package cron

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron/v3"
	"log/slog"
	"time"
	"webot/internal/bot"
	"webot/internal/cache"
	"webot/internal/constants"
	"webot/internal/context"
	"webot/internal/types"
)

func Start(ctx context.Context) {
	task := cron.New()
	_, _ = task.AddFunc(v2exSpec, v2ex(ctx).run)
	_, _ = task.AddFunc(aliveSpec, alive(ctx).run)
	task.Start()
	slog.Info("start cron", slog.String("v2ex", v2exSpec), slog.String("alive", aliveSpec))
}

func PushToAll(content string, key types.SettingsKey) error {
	self, err := bot.GetBot().GetCurrentUser()
	if err != nil {
		return err
	}
	friends, err := self.Friends()
	if err != nil {
		return err
	}
	var onFriends openwechat.Friends
	for _, f := range friends {
		if cache.GetSettings(f.AvatarID(), key) == constants.OFF {
			continue
		}
		onFriends = append(onFriends, f)
	}
	if len(onFriends) == 1 {
		onFriends = append(onFriends, onFriends[0])
	}
	if err = self.SendTextToFriends(content, time.Second, onFriends...); err != nil {
		return err
	}

	groups, err := self.Groups()
	if err != nil {
		return err
	}
	var onGroups openwechat.Groups
	for _, g := range groups {
		if cache.GetSettings(g.AvatarID(), key) == constants.OFF {
			continue
		}
		onGroups = append(onGroups, g)
	}

	if len(onGroups) == 1 {
		onGroups = append(onGroups, onGroups[0])
	}
	return self.SendTextToGroups(content, time.Second, onGroups...)
}
