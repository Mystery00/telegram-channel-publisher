package bot

import (
	"fmt"
	"slices"
	"strings"
	"telegram-channel-publisher/model"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	allowTypes = []string{"text_link", "hashtag", "pre", "bold", "italic", "underline", "strikethrough", "code"}
)

func ReadUrlByPhotoSizeList(photoList []tgbot.PhotoSize) (string, error) {
	maxFileId := ""
	maxFileSize := 0
	for _, p := range photoList {
		if p.FileSize > maxFileSize {
			maxFileId = p.FileID
			maxFileSize = p.FileSize
		}
	}
	return ReadUrl(maxFileId)
}

func ReadUrl(fileId string) (string, error) {
	file, err := bot.GetFile(tgbot.FileConfig{FileID: fileId})
	if err != nil {
		return "", err
	}
	link := file.Link(bot.Token)
	return link, nil
}

func DealEntities(list []tgbot.MessageEntity) []model.PostEntity {
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

func ReplaceApiEndpoint(post *model.Post) {
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
}
