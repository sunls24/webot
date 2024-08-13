package handler

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"log/slog"
	"time"
	"webot/internal/constants"
	"webot/internal/context"
	"webot/internal/entity"
	"webot/internal/types"
	"webot/pkg/openai"
)

type Handler struct {
	context.Context
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{ctx}
}

func (h *Handler) Handler(msg *openwechat.Message) {
	go h.handler(msg)
}

type params struct {
	atUser       string
	main, sender *openwechat.User
}

func (h *Handler) handler(msg *openwechat.Message) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("recover from panic", slog.Any("err", err))
		}
	}()

	if msg.IsFriendAdd() {
		onFriendAdd(msg)
		return
	}
	// TODO: 目前仅处理文本消息
	if !msg.IsText() || msg.IsSendBySelf() {
		return
	}

	sender, err := msg.Sender()
	if err != nil {
		slog.Error("get sender failed", slog.Any("err", err))
		return
	}

	//if !sender.IsPin() {
	//	return
	//}

	isGroup := sender.IsGroup()
	if isGroup && !msg.IsAt() {
		return
	}

	var ctx = params{
		main: sender,
	}
	if isGroup {
		atBot := []rune(fmt.Sprintf("@%s", msg.Owner().NickName))
		msgRune := []rune(msg.Content)
		if len(msgRune) <= len(atBot) {
			return
		}

		split := string(msgRune[len(atBot)])
		msg.Content = string(msgRune[len(atBot)+1:])
		if sender, err = msg.SenderInGroup(); err != nil {
			slog.Error("get sender in group failed", slog.Any("err", err))
			return
		}
		ctx.atUser = constants.AtUser(sender.NickName, split)
	}
	ctx.sender = sender

	slog.Info(fmt.Sprintf("%s: %s", ctx.main, msg.Content))
	switch {
	case isFunc(msg.Content, types.FuncHelp):
		h.onHelp(msg, ctx)
	case isFunc(msg.Content, types.FuncClear):
		h.onClear(msg, ctx)
	case isFunc(msg.Content, types.FuncV2ex):
		h.onToggle(msg, ctx, types.PushV2ex)
	default:
		h.onText(msg, ctx)
	}

	if err = msg.AsRead(); err != nil {
		slog.Error("msg.AsRead failed", slog.String("content", msg.Content), slog.Any("err", err))
	}
}

func onThink(msg *openwechat.Message, done <-chan struct{}, atUser string) {
	select {
	case <-done:
		return
	case <-time.After(time.Second * 2):
		_, err := msg.ReplyText(atUser + thinkText())
		if err != nil {
			slog.Error("reply think failed", slog.Any("err", err))
		}
	}
}

func (h *Handler) onText(msg *openwechat.Message, ctx params) {
	attr := slog.String("F", "onText")

	avatarID := ctx.main.AvatarID()
	messages, err := h.DB.GetMessagesByAvatarID(avatarID)
	if err != nil {
		slog.Error("get messages failed", attr, slog.Any("err", err))
		return
	}
	msgContent := msg.Content
	if ctx.atUser != "" {
		msgContent = fmt.Sprintf("%s: %s", ctx.sender.NickName, msgContent)
	}
	current := openai.Message{Role: openai.RUser, Content: msgContent}
	messages = BuildMessages(messages, current, ctx)

	done := make(chan struct{}, 1)
	go onThink(msg, done, ctx.atUser)
	start := time.Now()
	result, err := h.AI.Chat(h.Cfg.GetModel(false), messages)
	if err != nil {
		done <- struct{}{}
		slog.Error("AI chat failed", attr, slog.Any("err", err))
		_, err = msg.ReplyText(ctx.atUser + errorText(err))
		if err != nil {
			slog.Error("reply error failed", attr, slog.Any("err", err))
		}
		return
	}
	end := time.Now()
	done <- struct{}{}

	go func() {
		err = h.DB.InsertMessages(
			entity.NewMessage(avatarID, ctx.main.String(), &current),
			entity.NewMessage(avatarID, ctx.main.String(), result))
		if err != nil {
			slog.Error("insert messages failed", attr, slog.Any("err", err))
		}
		if len(messages) >= 20 {
			h.summarize(avatarID, ctx.main.String(), append(messages, *result))
		}
	}()

	_, err = msg.ReplyText(fmt.Sprintf(`%s
---------------------------
%s｜%.2fs`, ctx.atUser+result.Content, constants.HelpTip, end.Sub(start).Seconds()))
	if err != nil {
		slog.Error("reply text failed", slog.Any("err", err))
	}
}

func onFriendAdd(msg *openwechat.Message) {
	_ = msg.AsRead()
	f, err := msg.Agree()
	if err != nil {
		slog.Error("agree friend failed", slog.Any("err", err))
		return
	}
	slog.Info("agree friend add", slog.String("NickName", f.NickName))
	_, _ = f.SendText(helpText(f.User, f.User))
}
