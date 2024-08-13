package constants

import (
	"fmt"
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
	V2exTip = "V2EX 热帖推送: %s"
)

func Code(code string) string {
	if code == "" {
		return code
	}
	return fmt.Sprintf("```\n%s\n```", code)
}

const (
	OFF = "OFF"
	ON  = "ON"
)
