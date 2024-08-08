package database

import (
	"webot/internal/entity"
	"webot/pkg/openai"
)

func (db *DB) GetMessagesByAvatarID(avatarID string) ([]openai.Message, error) {
	var list = make([]entity.Message, 0)
	if err := db.Find(&list, "avatar_id = ?", avatarID).Error; err != nil {
		return nil, err
	}
	var result = make([]openai.Message, 0, len(list))
	for _, v := range list {
		result = append(result, v.OAI())
	}
	return result, nil
}

func (db *DB) InsertMessages(list ...entity.Message) error {
	return db.Create(&list).Error
}

func (db *DB) DeleteMessagesByAvatarID(avatarID string) error {
	return db.Delete(&entity.Message{}, "avatar_id = ?", avatarID).Error
}
