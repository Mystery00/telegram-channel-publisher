package publisher

import (
	"fmt"
	"net/url"
	"strings"
	"telegram-channel-publisher/model"
	"telegram-channel-publisher/publisher/halo"
	"time"
)

type HaloPublisher struct {
}

func (l *HaloPublisher) Publish(post model.Post) {
	attachments := make([]halo.Attachment, 0)
	if len(post.ImageList) > 0 {
		for _, v := range post.ImageList {
			u, mimeType := halo.DownloadAndUpload(v)
			attachments = append(attachments, halo.Attachment{
				FileType:     "PHOTO",
				FileMimeType: mimeType,
				FileURL:      u,
			})
		}
	}
	if len(post.VideoList) > 0 {
		for _, v := range post.VideoList {
			u, mimeType := halo.DownloadAndUpload(v)
			attachments = append(attachments, halo.Attachment{
				FileType:     "VIDEO",
				FileMimeType: mimeType,
				FileURL:      u,
			})
		}
	}
	//处理消息Entity
	blocks := make([]block, 0)
	for _, entity := range post.Entities {
		blocks = append(blocks, block{
			Type:       entity.Type,
			StartIndex: entity.Offset,
			EndIndex:   entity.Offset + entity.Length,
			Url:        entity.Url,
		})
	}
	content := strings.Builder{}
	originContent := post.Content
	tags := make([]string, 0)
	index := 0
	for i := range blocks {
		block := blocks[i]
		//将块之前的内容写入
		content.WriteString(parseContent(originContent, index, block.StartIndex, 0))
		//处理块
		switch block.Type {
		case "hashtag":
			s := parseContent(originContent, block.StartIndex, block.EndIndex, 1)
			tags = append(tags, s)
			esc := url.QueryEscape(s)
			r := fmt.Sprintf("<a class=\"tag\" href=\"?tag=%s\">%s</a>", esc, s)
			content.WriteString(r)
			break
		case "pre":
			s := parseContent(originContent, block.StartIndex, block.EndIndex, 0)
			content.WriteString(fmt.Sprintf("<pre><code>%s</code></pre>", s))
			break
		case "text_link":
			s := parseContent(originContent, block.StartIndex, block.EndIndex, 0)
			content.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a>", block.Url, s))
			break
		}
		index = block.EndIndex
	}
	//将剩余内容写入
	content.WriteString(parseContent(originContent, index, -1, 0))

	m := halo.Moment{
		Content:     content.String(),
		Attachments: attachments,
		ReleaseTime: time.Now().In(time.UTC),
		Tags:        tags,
	}
	halo.NewMoment(m)
}

type block struct {
	Type       string
	StartIndex int
	EndIndex   int
	Url        string
}

func parseContent(content string, start, end, offset int) string {
	r := []rune(content)
	if end == -1 {
		end = len(r)
	}
	return string(r[start+offset : end])
}
