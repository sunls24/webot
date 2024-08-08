package context

import (
	"webot/config"
	"webot/internal/database"
	"webot/pkg/openai"
)

type Context struct {
	Cfg *config.Config
	AI  *openai.OpenAI
	DB  *database.DB
}

func NewContext(cfg *config.Config) Context {
	db, err := database.Connect(cfg)
	if err != nil {
		panic(err)
	}
	return Context{
		DB:  db,
		Cfg: cfg,
		AI: openai.New(openai.Options{
			BaseURL: cfg.OpenAI.BaseURL,
			APIKey:  cfg.OpenAI.APIKey,
		}),
	}
}
