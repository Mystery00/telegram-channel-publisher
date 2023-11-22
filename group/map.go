package group

import (
	"sync"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RWMap struct {
	sync.RWMutex
	m map[string][]tgbot.Message
}

func NewRWMap() *RWMap {
	return &RWMap{
		m: make(map[string][]tgbot.Message),
	}
}

func (m *RWMap) Get(key string) ([]tgbot.Message, bool) {
	m.RLock()
	defer m.RUnlock()
	val, ok := m.m[key]
	return val, ok
}

func (m *RWMap) Add(key string, val tgbot.Message) {
	m.Lock()
	defer m.Unlock()
	m.m[key] = append(m.m[key], val)
}

func (m *RWMap) Group(f func(string, []tgbot.Message)) {
	m.Lock()
	arr := make(map[string][]tgbot.Message)
	for k, v := range m.m {
		arr[k] = v
	}
	m.m = make(map[string][]tgbot.Message)
	m.Unlock()
	for k, v := range arr {
		f(k, v)
	}
}
