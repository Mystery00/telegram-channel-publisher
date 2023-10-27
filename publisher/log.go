package publisher

import (
	"github.com/sirupsen/logrus"
	"telegram-channel-publisher/model"
)

type LogPublisher struct {
}

func (l *LogPublisher) Publish(post model.Post) {
	logrus.Debugf("publish post: %+v", post)
}
