package handler

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"webot/internal/cache"
	"webot/internal/constants"
	"webot/internal/types"
)

func thinkText() string {
	return fmt.Sprintf(constants.ThinkTip, openwechat.Emoji.Smart)
}

func errorText(err error) string {
	return fmt.Sprintf(constants.ErrorTip, openwechat.Emoji.Awkward, err)
}

func helpText(sender, main *openwechat.User) string {
	var avatarID = main.AvatarID()
	pushV2ex := cache.GetSettings(avatarID, types.PushV2ex)
	return fmt.Sprintf(`%s 你好，%s
%s 欢迎使用 %s:v%s
--------------------------
- help 查看帮助信息
- clear 清空聊天记录
- v2ex 热贴推送（%s）
- news 每日新闻推送（ON）`, openwechat.Emoji.LetMeSee, sender.NickName, openwechat.Emoji.Respect, constants.AppName, constants.Version, pushV2ex)
}

func clearText() string {
	return fmt.Sprintf(constants.ClearTip, openwechat.Emoji.Doge)
}

func v2exText(s string) string {
	return fmt.Sprintf(constants.V2exTip, s)
}
