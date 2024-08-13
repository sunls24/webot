package database

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"webot/config"
	"webot/internal/entity"
)

type DB struct {
	*gorm.DB
}

func Connect(cfg *config.Config) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	if err = migrate(db); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return &DB{db}, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.Message{}, &entity.V2ex{}, &entity.Settings{})
}
