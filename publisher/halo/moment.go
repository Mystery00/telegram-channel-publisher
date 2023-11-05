package halo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

const templateDir = "D:\\GolandProjects\\telegram-channel-publisher\\templates"

func NewMoment(m Moment) {
	medium, err := json.Marshal(m.Attachments)
	if err != nil {
		panic(err)
	}
	html := convertHtml(m.Content)
	tags, err := json.Marshal(m.Tags)
	if err != nil {
		panic(err)
	}
	tt := momentTpl{
		Content:     html,
		Html:        html,
		Medium:      string(medium),
		ReleaseTime: m.ReleaseTime,
		Tags:        string(tags),
	}
	tplPath := fmt.Sprintf("%s/halo-moment.tpl", templateDir)
	t := template.Must(template.ParseGlob(tplPath))
	var tmplBytes bytes.Buffer
	err = t.Execute(&tmplBytes, tt)
	if err != nil {
		panic(err)
	}
	b := tmplBytes.Bytes()
	logrus.Debugf("publish moment: %s", string(b))
	_, err = request("apis/api.plugin.halo.run/v1alpha1/plugins/PluginMoments/moments", http.MethodPost, b)
	if err != nil {
		logrus.Errorf("publish error: %v", err)
	}
	logrus.Infof("publish moment success")
}

func convertHtml(content string) string {
	c := strings.Builder{}
	for _, s := range strings.Split(content, "\n") {
		c.WriteString(fmt.Sprintf("<p>%s</p>", s))
	}
	return strings.ReplaceAll(c.String(), "\"", "\\\"")
}

type momentTpl struct {
	Content     string
	Html        string
	Medium      string
	ReleaseTime time.Time
	Tags        string
}

type Moment struct {
	Content     string
	Tags        []string
	Attachments []Attachment
	ReleaseTime time.Time
}

type Attachment struct {
	FileType     string `json:"type"`
	FileMimeType string `json:"originType"`
	FileURL      string `json:"url"`
}
