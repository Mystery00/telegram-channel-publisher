package publisher

import (
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
			url, mimeType := halo.DownloadAndUpload(v)
			attachments = append(attachments, halo.Attachment{
				FileType:     "PHOTO",
				FileMimeType: mimeType,
				FileURL:      url,
			})
		}
	}
	if len(post.VideoList) > 0 {
		for _, v := range post.VideoList {
			url, mimeType := halo.DownloadAndUpload(v)
			attachments = append(attachments, halo.Attachment{
				FileType:     "VIDEO",
				FileMimeType: mimeType,
				FileURL:      url,
			})
		}
	}
	m := halo.Moment{
		Content:     post.Content,
		Attachments: attachments,
		ReleaseTime: time.Now().In(time.UTC),
	}
	halo.NewMoment(m)
}
