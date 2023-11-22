package group

import (
	"telegram-channel-publisher/bot"
	"telegram-channel-publisher/config"
	"telegram-channel-publisher/model"
	"telegram-channel-publisher/publisher"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	delayTimeS = viper.GetInt64(config.MediaDelay)

	prepareHandleMap = NewRWMap()
)

func AddMediaMessage(m tgbot.Message) {
	//计算执行时间
	executeTimeMs := time.Now().Add(time.Duration(delayTimeS) * time.Second).UnixMilli()
	prepareHandleMap.Add(m.MediaGroupID, m, executeTimeMs)
}

func InitMediaHandleCron() {
	cronjob := cron.New(cron.WithSeconds())
	_, err := cronjob.AddFunc("*/5 * * * * *", doHandle)
	if err != nil {
		logrus.Fatalf("register cron job error: %v", err)
	}
	logrus.Debugf("register cron job success")
	cronjob.Start()
}

func doHandle() {
	prepareHandleMap.Group(func(key string, m MediaGroup) {
		if len(m.Messages) == 0 {
			return
		}
		logrus.Debugf("handle media group: %s", key)
		//获取第一个消息
		msg := m.Messages[0]
		post := model.Post{}
		post.Content = msg.Caption
		post.Entities = bot.DealEntities(msg.CaptionEntities)
		imageList := make([]string, 0)
		for _, message := range m.Messages {
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
