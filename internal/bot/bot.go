package bot

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"log/slog"
	"time"
	"webot/internal/context"
	"webot/internal/handler"
)

var bot *openwechat.Bot

func GetBot() *openwechat.Bot {
	return bot
}

const (
	retry = time.Minute * 10
)

func Block(ctx context.Context) {

	for {
		if err := syncBot(ctx); err != nil {
			slog.Error("trying login again", slog.Duration("retry", retry), slog.Any("err", err))
		} else {
			break
		}
		<-time.After(retry)
	}
}

func syncBot(ctx context.Context) error {
	bot = openwechat.DefaultBot(ctx.Cfg.GetMode())
	bot.MessageHandler = handler.NewHandler(ctx).Handler
	bot.SyncCheckCallback = syncCheckCallback
	bot.UUIDCallback = printlnQrcodeUrl
	bot.LoginCallBack = loginCallback
	bot.ScanCallBack = scanCallback

	slog.Info("start login", slog.String("mode", string(ctx.Cfg.Mode)))
	var storage = ctx.Cfg.NewStorage()
	if err := bot.HotLogin(storage); err != nil {
		_ = storage.Close()
		slog.Warn("hot login failed, try push login", slog.Any("err", err))
		storage = ctx.Cfg.NewStorage()
		if err = bot.PushLogin(storage, openwechat.NewRetryLoginOption()); err != nil {
			_ = storage.Close()
			return fmt.Errorf("push login: %w", err)
		}
	}
	defer storage.Close()
	slog.Info("login complete, start sync messages")

	return bot.Block()
}

func printlnQrcodeUrl(uuid string) {
	qrcodeUrl := openwechat.GetQrcodeUrl(uuid)
	slog.Info("访问下面网址扫描二维码登录")
	slog.Info(qrcodeUrl)
}

func scanCallback(resp openwechat.CheckLoginResponse) {
	slog.Info("扫码成功，请在手机上确认登录")
}

func loginCallback(resp openwechat.CheckLoginResponse) {
	slog.Info("登录成功")
}

func syncCheckCallback(resp openwechat.SyncCheckResponse) {
	if resp.Success() && resp.NorMal() {
		return
	}
	if resp.HasNewMessage() {
		return
	}
	slog.Info("sync check", slog.String("RetCode", resp.RetCode), slog.String("Selector", string(resp.Selector)))
}
