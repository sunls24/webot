package main

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"webot/config"
	"webot/internal/app"
	"webot/internal/constants"
)

func main() {
	f := flag.NewFlagSet(constants.AppName, flag.ExitOnError)
	f.Usage = cleanenv.FUsage(f.Output(), &config.Config{}, nil, f.Usage)

	var cfgPath string
	f.StringVar(&cfgPath, "cfg", "", "path to config file")
	_ = f.Parse(os.Args[1:])

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		panic(err)
	}

	app.Run(cfg)
}
