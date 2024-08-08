package entity

import (
	"gorm.io/gorm"
	"webot/pkg/openai"
)

type Message struct {
	gorm.Model

	AvatarID string `gorm:"index"`
	NickName string
	Role     string
	Content  string
}

func (m *Message) OAI() openai.Message {
	return openai.Message{
		Role:    openai.Role(m.Role),
		Content: m.Content,
	}
}

func NewMessage(avatarID, nickName string, msg *openai.Message) Message {
	return Message{
		AvatarID: avatarID,
		NickName: nickName,
		Role:     string(msg.Role),
		Content:  msg.Content,
	}
}
