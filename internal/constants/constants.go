package constants

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
)

const (
	AppName = "webot"
	Version = "0.1.0"
)

const Welcome = `                __          __ 
 _      _____  / /_  ____  / /_
| | /| / / _ \/ __ \/ __ \/ __/
| |/ |/ /  __/ /_/ / /_/ / /_  
|__/|__/\___/_.___/\____/\__/
`

func AtUser(nickName, split string) string {
	return fmt.Sprintf("@%s%s\n", nickName, split)
}

const (
	HelpTip  = "输入 help 查看帮助"
	ThinkTip = "%s 正在思考中，请稍等..."

	ClearTip = "%s 清空了，已开启新的对话！"
	ErrorTip = `%s 哎呀，出错了！
------------------
%v`
)

func HelpText(sender *openwechat.User) string {
	return fmt.Sprintf(`%s 你好，%s
%s 欢迎使用 %s:v%s
--------------------------
- help 查看帮助信息
- clear 清空聊天记录
- v2ex 热贴推送（关闭）
- news 每日新闻推送（关闭）`, openwechat.Emoji.LetMeSee, sender.NickName, openwechat.Emoji.Respect, AppName, Version)
}

func ThinkText() string {
	return fmt.Sprintf(ThinkTip, openwechat.Emoji.Smart)
}

func ClearText() string {
	return fmt.Sprintf(ClearTip, openwechat.Emoji.Doge)
}

func ErrorText(err error) string {
	return fmt.Sprintf(ErrorTip, openwechat.Emoji.Awkward, err)
}
