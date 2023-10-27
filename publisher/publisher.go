package publisher

import (
	"github.com/spf13/viper"
	"telegram-channel-publisher/config"
	"telegram-channel-publisher/model"
)

type Publisher interface {
	// Publish 发布消息
	Publish(post model.Post)
}

var publisher Publisher

func InitPublisher() {
	pubType := viper.GetString(config.PublisherType)
	switch pubType {
	default:
		publisher = &LogPublisher{}
		break
	}
}

func Pub(post model.Post) {
	publisher.Publish(post)
}
