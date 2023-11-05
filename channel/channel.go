package channel

import (
	"fmt"
	"slices"
	"strings"
	"telegram-channel-publisher/config"
	"telegram-channel-publisher/model"
	"telegram-channel-publisher/publisher"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	apiEndpoint = viper.GetString(config.ApiEndpoint)
	channel     = viper.GetInt64(config.Channel)
	filterTags  = viper.GetStringSlice(config.FilterTags)

	allowTypes = []string{"text_link", "hashtag", "pre", "bold", "italic", "underline", "strikethrough", "code"}
)

func HandleUpdate(bot *tgbot.BotAPI, inCh <-chan tgbot.Update) {
	go func() {
		for ch := range inCh {
			if ch.ChannelPost == nil {
				//不是频道消息，跳过
				logrus.Debugf("not channel post, skip")
				continue
			}
			logrus.Debugf("receive channel post from: %d", ch.ChannelPost.Chat.ID)
			if channel != 0 && ch.ChannelPost.Chat.ID != channel {
				//不是指定频道的消息，跳过
				logrus.Debugf("not channel [%d] post, skip", channel)
				continue
			}
			if len(filterTags) > 0 {
				//有过滤标签，检查是否包含
				skip := false
				for _, tag := range filterTags {
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
			msg := ch.ChannelPost
			post := model.Post{
				Sender: msg.SenderChat.UserName,
			}
			if msg.Animation != nil {
				//GIF
				post.Content = msg.Caption
				post.Entities = dealEntities(post.Content, msg.CaptionEntities)
				url, err := readUrl(bot, msg.Animation.FileID)
				if err != nil {
					logrus.Errorf("read url error: %v", err)
					continue
				}
				post.WithVideo(url)
			} else if msg.Sticker != nil {
				//贴纸包
				url, err := readUrl(bot, msg.Sticker.FileID)
				if err != nil {
					logrus.Errorf("read url error: %v", err)
					continue
				}
				post.WithImage(url)
			} else if msg.Photo != nil {
				//带图片的消息
				post.Content = msg.Caption
				post.Entities = dealEntities(post.Content, msg.CaptionEntities)
				if len(msg.Photo) == 0 {
					continue
				}
				urlList, err := readUrlByPhotoSizeList(bot, msg.Photo)
				if err != nil {
					logrus.Errorf("read url list error: %v", err)
					continue
				}
				post.WithImage(urlList)
			} else if msg.Text != "" {
				//纯文本消息
				post.Content = msg.Text
				post.Entities = dealEntities(post.Content, msg.Entities)
			} else {
				//其他类型消息，跳过
				continue
			}
			//替换地址
			if apiEndpoint != "" {
				if len(post.ImageList) > 0 {
					replaceHost := strings.TrimSuffix(apiEndpoint, "/bot%s/%s")
					for i := range post.ImageList {
						path := strings.TrimPrefix(post.ImageList[i], "https://api.telegram.org/")
						post.ImageList[i] = fmt.Sprintf("%s/%s", replaceHost, path)
					}
				}
				if len(post.VideoList) > 0 {
					replaceHost := strings.TrimSuffix(apiEndpoint, "/bot%s/%s")
					for i := range post.VideoList {
						path := strings.TrimPrefix(post.VideoList[i], "https://api.telegram.org/")
						post.VideoList[i] = fmt.Sprintf("%s/%s", replaceHost, path)
					}
				}
			}
			go publisher.Pub(post)
		}
	}()
}

func readUrlByPhotoSizeList(bot *tgbot.BotAPI, photoList []tgbot.PhotoSize) (string, error) {
	maxFileId := ""
	maxFileSize := 0
	for _, p := range photoList {
		if p.FileSize > maxFileSize {
			maxFileId = p.FileID
			maxFileSize = p.FileSize
		}
	}
	return readUrl(bot, maxFileId)
}

func readUrl(bot *tgbot.BotAPI, fileId string) (string, error) {
	file, err := bot.GetFile(tgbot.FileConfig{FileID: fileId})
	if err != nil {
		return "", err
	}
	link := file.Link(bot.Token)
	return link, nil
}

func dealEntities(text string, list []tgbot.MessageEntity) []model.PostEntity {
	entities := make([]model.PostEntity, 0)
	for _, entity := range list {
		if !slices.Contains(allowTypes, entity.Type) {
			continue
		}
		entities = append(entities, model.PostEntity{
			Type:   entity.Type,
			Offset: entity.Offset,
			Length: entity.Length,
			Url:    entity.URL,
		})
	}
	return entities
}
