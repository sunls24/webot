package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"webot/internal/types"
)

type (
	Config struct {
		Mode    types.LoginMode `yaml:"mode" env:"MODE" env-default:"desktop"`
		Storage string          `yaml:"storage" env:"STORAGE" env-default:"webot.json"`

		OpenAI OpenAI `yaml:"openai"`
		DB     DB     `yaml:"db"`
	}

	OpenAI struct {
		BaseURL   string `yaml:"base_url" env:"OPENAI_BASE_URL" env-default:"https://api.openai.com"`
		APIKey    string `yaml:"api_key" env:"OPENAI_API_KEY" env-required:"true"`
		Model     string `yaml:"model" env:"OPENAI_MODEL" env-default:"gpt-3.5-turbo"`
		SlowModel string `yaml:"slow_model" env:"OPENAI_SLOW_MODEL" env-default:"gpt-4"`
	}

	DB struct {
		DSN string `yaml:"dsn" env:"DB_DSN" env-default:"webot.db"`
	}
)

func NewConfig(cfgPath string) (*Config, error) {
	cfg := &Config{}
	if cfgPath != "" {
		if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
			return nil, err
		}
	}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
