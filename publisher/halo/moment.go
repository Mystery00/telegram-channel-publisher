package halo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

var linkReg = regexp.MustCompile(`(http(s)://([^ \n]+))`)

const templateDir = "/app/templates"

func NewMoment(m Moment) {
	medium, err := json.Marshal(m.Attachments)
	if err != nil {
		panic(err)
	}
	tt := momentTpl{
		Content:     convertStrToHtml(m.Content),
		Html:        convertStrToHtml(m.Content),
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
	_, err = request("apis/api.plugin.halo.run/v1alpha1/plugins/PluginMoments/moments", http.MethodPost, b)
	if err != nil {
		logrus.Errorf("publish error: %v", err)
	}
	logrus.Infof("publish moment success")
}

func convertStrToHtml(str string) string {
	sb := strings.Builder{}
	for _, s := range strings.Split(str, "\n") {
		if s == "" {
			continue
		}
		//convert url to a tag
		ss := linkReg.ReplaceAllString(s, "<a href=\"$1\">$1</a>")
		sb.WriteString(fmt.Sprintf("<p>%s</p>", ss))
	}
	return strings.ReplaceAll(sb.String(), "\"", "\\\"")
}

type momentTpl struct {
	Content     string
	Html        string
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
