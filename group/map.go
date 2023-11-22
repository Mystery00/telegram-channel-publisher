package group

import (
	"sync"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RWMap struct {
	sync.RWMutex
	m map[string]MediaGroup
}

func NewRWMap() *RWMap {
	return &RWMap{
		m: make(map[string]MediaGroup),
	}
}

func (m *RWMap) Get(key string) (MediaGroup, bool) {
	m.RLock()
	defer m.RUnlock()
	val, ok := m.m[key]
	return val, ok
}

func (m *RWMap) Add(key string, msg tgbot.Message, executeTimeMs int64) {
	m.Lock()
	defer m.Unlock()
	mediaGroup, exist := m.m[key]
	var mg *MediaGroup
	if !exist {
		mg = NewMediaGroup(executeTimeMs)
	} else {
		mg = &mediaGroup
	}
	mg.AddMessage(executeTimeMs, msg)
	m.m[key] = *mg
}

func (m *RWMap) Group(f func(string, MediaGroup)) {
	m.Lock()
	arr := make(map[string]MediaGroup)
	notExpired := make(map[string]MediaGroup)
	now := time.Now().UnixMilli()
	for k, v := range m.m {
		if v.ExpectExecuteTimeMs > now {
			notExpired[k] = v
		} else {
			arr[k] = v
		}
	}
	m.m = notExpired
	m.Unlock()
	for k, v := range arr {
		f(k, v)
	}
}
