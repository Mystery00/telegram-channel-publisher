package model

type Post struct {
	Content   string
	Sender    string
	ImageList []string
	VideoList []string
}

func (p *Post) WithImage(imageUrl string) *Post {
	p.ImageList = make([]string, 0)
	p.ImageList = append(p.ImageList, imageUrl)
	return p
}

func (p *Post) WithImages(imageList []string) *Post {
	p.ImageList = imageList
	return p
}

func (p *Post) WithVideo(videoUrl string) *Post {
	p.VideoList = make([]string, 0)
	p.VideoList = append(p.VideoList, videoUrl)
	return p
}

func (p *Post) WithVideos(videoList []string) *Post {
	p.VideoList = videoList
	return p
}
