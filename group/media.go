package group

import tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MediaGroup struct {
	ExpectExecuteTimeMs int64
	Messages            []tgbot.Message
}

func NewMediaGroup(executeTimeMs int64) *MediaGroup {
	return &MediaGroup{
		ExpectExecuteTimeMs: executeTimeMs,
		Messages:            make([]tgbot.Message, 0),
	}
}

func (m *MediaGroup) AddMessage(newExecuteTimeMs int64, message tgbot.Message) {
	m.ExpectExecuteTimeMs = newExecuteTimeMs
	m.Messages = append(m.Messages, message)
}
