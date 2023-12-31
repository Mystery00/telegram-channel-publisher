package channel

import (
	"fmt"
	"strings"
	"telegram-channel-publisher/bot"
	"telegram-channel-publisher/config"
	"telegram-channel-publisher/group"
	"telegram-channel-publisher/model"
	"telegram-channel-publisher/publisher"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	delayTimeS = viper.GetInt64(config.MediaDelay)

	channelEnable = viper.GetBool(config.ChannelEnable)
	channelId     = viper.GetInt64(config.ChannelId)
	channelFilter = viper.GetStringSlice(config.ChannelFilter)

	privateEnable = viper.GetBool(config.PrivateEnable)
	privateSender = viper.GetInt64(config.PrivateSender)
)

func HandleUpdate(inCh <-chan tgbot.Update) {
	go func() {
		for ch := range inCh {
			if ch.ChannelPost != nil {
				//频道消息
				logrus.Debugf("receive channel post from: %d", ch.ChannelPost.Chat.ID)
				if !channelEnable {
					logrus.Debugf("channel post disable, skip")
					continue
				}
				if channelId != 0 && ch.ChannelPost.Chat.ID != channelId {
					//不是指定频道的消息，跳过
					logrus.Debugf("not channel [%d] post, skip", channelId)
					continue
				}
				if len(channelFilter) > 0 {
					//有过滤标签，检查是否包含
					skip := false
					for _, tag := range channelFilter {
						f := fmt.Sprintf("#%s", tag)
						if strings.Contains(ch.ChannelPost.Text, f) {
							//包含指定标签，跳过
							logrus.Infof("contain tag [%s] post, skip", tag)
							skip = true
							continue
						}
					}
					if skip {
						continue
					}
				}
				handleMessage(ch.ChannelPost, false)
			} else if ch.Message != nil && ch.Message.Chat.IsPrivate() {
				//私聊消息
				logrus.Debugf("receive private message from: %d", ch.Message.Chat.ID)
				if !privateEnable {
					logrus.Debugf("private message disable, skip")
					continue
				}
				if privateSender != 0 && ch.Message.Chat.ID != privateSender {
					logrus.Debugf("not private sender [%d] message, skip", privateSender)
					continue
				}
				handleMessage(ch.Message, true)
			}
		}
	}()
}

func handleMessage(msg *tgbot.Message, isPrivate bool) {
	post := model.Post{}
	if msg.Animation != nil {
		//GIF
		post.Content = msg.Caption
		post.Entities = bot.DealEntities(msg.CaptionEntities)
		url, err := bot.ReadUrl(msg.Animation.FileID)
		if err != nil {
			logrus.Errorf("read url error: %v", err)
			return
		}
		post.WithVideo(url)
	} else if msg.Sticker != nil {
		//贴纸包
		url, err := bot.ReadUrl(msg.Sticker.FileID)
		if err != nil {
			logrus.Errorf("read url error: %v", err)
			return
		}
		post.WithImage(url)
	} else if msg.Photo != nil {
		//带图片的消息
		if delayTimeS != 0 && msg.MediaGroupID != "" {
			//图片组
			group.AddMediaMessage(*msg)
			return
		} else {
			post.Content = msg.Caption
			post.Entities = bot.DealEntities(msg.CaptionEntities)
			if len(msg.Photo) == 0 {
				return
			}
			urlList, err := bot.ReadUrlByPhotoSizeList(msg.Photo)
			if err != nil {
				logrus.Errorf("read url list error: %v", err)
				return
			}
			post.WithImage(urlList)
		}
	} else if msg.Text != "" {
		//纯文本消息
		post.Content = msg.Text
		post.Entities = bot.DealEntities(msg.Entities)
	} else {
		//其他类型消息，跳过
		return
	}
	//替换地址
	bot.ReplaceApiEndpoint(&post)
	//设置消息ID
	post.MessageId = msg.MessageID
	post.ChatId = msg.Chat.ID
	post.IsPrivate = isPrivate
	go publisher.Pub(post)
}
