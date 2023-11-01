package halo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

const templateDir = "/app/templates"

func NewMoment(m Moment) {
	medium, err := json.Marshal(m.Attachments)
	if err != nil {
		panic(err)
	}
	tt := momentTpl{
		Content:     m.Content,
		Medium:      string(medium),
		ReleaseTime: m.ReleaseTime,
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
	_, err = request("/apis/api.plugin.halo.run/v1alpha1/plugins/PluginMoments/moments", http.MethodPost, b)
	if err != nil {
		logrus.Errorf("publish error: %v", err)
	}
	logrus.Infof("publish moment success")
}

type momentTpl struct {
	Content     string
	Medium      string
	ReleaseTime time.Time
}

type Moment struct {
	Content     string
	Attachments []Attachment
	ReleaseTime time.Time
}

type Attachment struct {
	FileType     string `json:"type"`
	FileMimeType string `json:"originType"`
	FileURL      string `json:"url"`
}
