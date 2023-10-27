package channel

import (
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"telegram-channel-publisher/model"
	"telegram-channel-publisher/publisher"
)

func HandleUpdate(bot *tgbot.BotAPI, inCh <-chan tgbot.Update) {
	go func() {
		for ch := range inCh {
			if ch.ChannelPost == nil {
				//不是频道消息，跳过
				continue
			}
			msg := ch.ChannelPost
			post := model.Post{
				Sender: msg.SenderChat.UserName,
			}
			if msg.Animation != nil {
				//GIF
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
			} else {
				//其他类型消息，跳过
				continue
			}
			publisher.Pub(post)
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
