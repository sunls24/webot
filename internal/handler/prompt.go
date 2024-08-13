package handler

import (
	"fmt"
	"log/slog"
	"time"
	"webot/internal/entity"
	"webot/pkg/openai"
)

const (
	system = `You are a WeChat chatbot, your duty is to provide accurate and rigorous answers. 
Please follow these requirements:
1. Please keep your reply concise, avoid lengthy explanations, preferably use Chinese
2. Do not use Markdown format, use plain text format
Current time is: %s`
	summarize = "简要总结一下对话内容，用作后续的上下文提示 prompt，控制在 200 字以内"
)

func systemPrompt(ctx params) openai.Message {
	var group string
	if ctx.atUser != "" {
		group = fmt.Sprintf("\nCurrently running in group '%s', message will begin with member's nickname", ctx.main.NickName)
	}
	return openai.Message{
		Role:    openai.RSystem,
		Content: fmt.Sprintf(system, time.Now().Format(time.DateTime)) + group,
	}
}
func BuildMessages(list []openai.Message, current openai.Message, ctx params) []openai.Message {
	var result = make([]openai.Message, 0, len(list)+2)
	result = append(result, systemPrompt(ctx))
	if len(list) > 0 {
		result = append(result, list...)
	}
	result = append(result, current)
	return result
}

func (h *Handler) summarize(avatarID, nickName string, messages []openai.Message) {
	attr := slog.String("F", "summarize")

	messages = append(messages, openai.Message{
		Role:    openai.RUser,
		Content: summarize,
	})

	result, err := h.AI.Chat(h.Cfg.GetModel(true), messages)
	if err != nil {
		slog.Error("AI chat failed", attr, slog.Any("err", err))
		return
	}

	err = h.DB.DeleteMessagesByAvatarID(avatarID)
	if err != nil {
		slog.Error("delete messages failed", attr, slog.Any("err", err))
	}

	err = h.DB.InsertMessages(entity.NewMessage(avatarID, nickName, result))
	if err != nil {
		slog.Error("insert messages failed", attr, slog.Any("err", err))
	}
}
