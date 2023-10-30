package publisher

import (
	"telegram-channel-publisher/model"
	"telegram-channel-publisher/publisher/halo"
	"time"
)

type HaloPublisher struct {
}

func (l *HaloPublisher) Publish(post model.Post) {
	imageUrl := ""
	mimeType := ""
	if len(post.ImageList) > 0 {
		imageUrl, mimeType = halo.DownloadAndUpload(post.ImageList[0])
	}
	m := halo.Moment{
		Content:       post.Content,
		ImageMimeType: mimeType,
		ImageURL:      imageUrl,
		ReleaseTime:   time.Now().In(time.UTC),
	}
	halo.NewMoment(m)
}
