package group

import (
	"telegram-channel-publisher/bot"
	"telegram-channel-publisher/model"
	"telegram-channel-publisher/publisher"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var prepareHandleMap = NewRWMap()

func AddMediaMessage(m tgbot.Message) {
	prepareHandleMap.Add(m.MediaGroupID, m)
}

func InitMediaHandleCron() {
	cronjob := cron.New(cron.WithSeconds())
	_, err := cronjob.AddFunc("*/5 * * * * *", doHandle)
	if err != nil {
		logrus.Errorf("register cron job error: %v", err)
		return
	}
	logrus.Debugf("register cron job success")
	cronjob.Start()
}

func doHandle() {
	prepareHandleMap.Group(func(key string, val []tgbot.Message) {
		if len(val) == 0 {
			return
		}
		logrus.Debugf("handle media group: %s", key)
		//获取第一个消息
		msg := val[0]
		post := model.Post{}
		post.Content = msg.Caption
		post.Entities = bot.DealEntities(msg.CaptionEntities)
		imageList := make([]string, 0)
		for _, message := range val {
			if len(message.Photo) == 0 {
				continue
			}
			u, err := bot.ReadUrlByPhotoSizeList(message.Photo)
			if err != nil {
				logrus.Errorf("read url list error: %v", err)
				return
			}
			imageList = append(imageList, u)
		}
		post.WithImages(imageList)
		//替换地址
		bot.ReplaceApiEndpoint(&post)
		go publisher.Pub(post)
	})
}
