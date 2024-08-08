package app

import (
	"fmt"
	"log/slog"
	"webot/config"
	"webot/internal/bot"
	"webot/internal/constants"
	"webot/internal/context"
)

func Run(cfg *config.Config) {
	//goland:noinspection GoPrintFunctions
	fmt.Println(constants.Welcome)
	fmt.Printf("======== %s:v%s ========\n", constants.AppName, constants.Version)

	slog.Info("login mode: " + string(cfg.Mode))
	ctx := context.NewContext(cfg)
	bot.Block(ctx)
}
