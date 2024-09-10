package cron

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"webot/internal/bot"
	"webot/internal/constants"
	"webot/internal/context"
	"webot/internal/entity"
	"webot/internal/types"
	"webot/pkg/client"
	"webot/pkg/openai"

	"github.com/tidwall/gjson"
)

const (
	v2exSpec       = "30 10-18/3 * * *"
	v2exHotAPI     = "https://www.v2ex.com/api/topics/hot.json"
	v2exRepliesAPI = "https://www.v2ex.com/api/replies/show.json?topic_id="
)

type v2ex context.Context

func (ctx v2ex) run() {
	slog.Info("cron v2ex run")
	attr := slog.String("F", "v2ex.run")
	if !bot.GetBot().Alive() {
		slog.Warn("bot is not alive", attr)
		return
	}

	req, _ := http.NewRequest(http.MethodGet, v2exHotAPI, nil)
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("get v2ex hot api failed", slog.Any("err", err))
		return
	}
	summarizeList := make([]string, 0)
	for _, v := range gjson.ParseBytes(resp).Array() {
		id := v.Get("id").String()
		if exist := ctx.DB.V2exExist(id); exist {
			continue
		}

		title := v.Get("title").String()
		content := v.Get("content").String()
		member := v.Get("member.username").String()
		url := v.Get("url").String()

		summarize, err := ctx.v2exSummarize(id, member, title, content, url)
		if err != nil {
			slog.Error("v2ex summarize failed", slog.Any("err", err))
			continue
		}

		slog.Info("v2ex post", slog.Any("title", title))
		summarizeList = append(summarizeList, summarize)
		ctx.DB.Create(&entity.V2ex{
			ID:    id,
			Title: title,
		})
	}
	if len(summarizeList) == 0 {
		return
	}
	if err = PushToAll(v2exHeader+strings.Join(summarizeList, v2exSplit)+v2exSplit+v2exFooter, types.PushV2ex); err != nil {
		slog.Error("push v2ex post failed", slog.Any("err", err))
	}
}

func (ctx v2ex) v2exSummarize(id, member, title, content, url string) (string, error) {
	req, _ := http.NewRequest(http.MethodGet, v2exRepliesAPI+id, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	replies := gjson.ParseBytes(resp).Array()
	repliesStr := make([]string, 0, len(replies))
	for _, r := range replies {
		repliesStr = append(repliesStr, fmt.Sprintf("- %s:\n%s",
			r.Get("member.username").String(),
			r.Get("content").String()))
	}

	if content == "" {
		content = title
	}
	messages := []openai.Message{
		{Role: openai.RSystem, Content: v2exSystem},
		{Role: openai.RUser, Content: fmt.Sprintf(`帖子内容如下:
标题：%s
作者：%s
内容：%s
回复列表：
%s`, title, member, constants.Code(content), constants.Code(strings.Join(repliesStr, "\n")))},
	}

	result, err := ctx.AI.Chat(ctx.Cfg.GetModel(true), messages)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`主题：%s
AI 点评：%s
链接：%s`, title, result.Content, url), nil
}

const (
	v2exSystem = "你擅长点评 V2EX 论坛热门帖子，从主题内容和回复中提取关键信息，并使用 60 个汉字以内的简洁描述。"
	v2exSplit  = "\n-------------------\n"

	v2exHeader = `V2EX 热帖推送
---------------
`
	v2exFooter = "输入 v2ex 关闭推送"
)
