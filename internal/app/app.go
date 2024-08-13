package app

import (
	"fmt"
	"webot/config"
	"webot/internal/bot"
	"webot/internal/cache"
	"webot/internal/constants"
	"webot/internal/context"
	"webot/internal/cron"
)

func Run(cfg *config.Config) {
	//goland:noinspection GoPrintFunctions
	fmt.Println(constants.Welcome)
	fmt.Printf("======== %s:v%s ========\n", constants.AppName, constants.Version)

	ctx := context.NewContext(cfg)
	cache.InitSettings(ctx.DB)
	cron.Start(ctx)
	bot.Block(ctx)
}
