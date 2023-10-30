package halo

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"telegram-channel-publisher/config"
	"time"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var (
	group  = viper.GetString(config.HaloImageGroup)
	policy = viper.GetString(config.HaloImagePolicy)
)

func DownloadAndUpload(imageUrl string) (string, string) {
	//下载文件到本地
	response, err := http.Get(imageUrl)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		panic(fmt.Errorf("download image error: %d", response.StatusCode))
	}
	//截取文件名
	i := strings.LastIndex(imageUrl, ".")
	extName := imageUrl[i+1:]
	fileName := fmt.Sprintf("moments-%d.%s", time.Now().UnixMilli(), extName)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, response.Body)
	if err != nil {
		panic(err)
	}
	_ = writer.WriteField("groupName", group)
	_ = writer.WriteField("policyName", policy)
	err = writer.Close()
	if err != nil {
		panic(err)
	}
	resp, err := requestWithContentType("/apis/api.console.halo.run/v1alpha1/attachments/upload", http.MethodPost, writer.FormDataContentType(), body)
	if err != nil {
		panic(err)
	}
	attachmentName := gjson.GetBytes(resp, "metadata.name").String()
	//查询附件详情
	index := 0
	for {
		if index > 100 {
			//超时
			panic(fmt.Errorf("wait for attachment timeout"))
		}
		time.Sleep(time.Millisecond * 100)
		attachmentInfo, err := request(fmt.Sprintf("/apis/storage.halo.run/v1alpha1/attachments/%s", attachmentName), http.MethodGet, nil)
		if err != nil {
			panic(err)
		}
		permalink := gjson.GetBytes(attachmentInfo, "status.permalink").String()
		if permalink != "" {
			return permalink, detectImageMimeType(imageUrl)
		}
		index++
	}
}

func detectImageMimeType(imageUrl string) string {
	if strings.HasSuffix(imageUrl, ".jpg") || strings.HasSuffix(imageUrl, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(imageUrl, ".png") {
		return "image/png"
	} else if strings.HasSuffix(imageUrl, ".webp") {
		return "image/webp"
	} else {
		panic(fmt.Errorf("unsupported image type for: %s", imageUrl))
	}
}
