package halo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

var linkReg = regexp.MustCompile(`(http(s)://([^ \n]+))`)
var tagReg = regexp.MustCompile(`#([^ \n]+)`)

const templateDir = "templates"

func NewMoment(m Moment) {
	medium, err := json.Marshal(m.Attachments)
	if err != nil {
		panic(err)
	}
	html, tagList := convertStrToHtml(m.Content)
	tags, err := json.Marshal(tagList)
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

func convertStrToHtml(str string) (string, []string) {
	sb := strings.Builder{}
	tags := make([]string, 0)
	for _, s := range strings.Split(str, "\n") {
		if s == "" {
			continue
		}
		if strings.HasPrefix(s, "#") {
			//判定为Tag
			allTagArray := tagReg.FindAllStringSubmatch(s, -1)
			if len(allTagArray) > 0 {
				sb.WriteString("<p>")
				for _, tag := range allTagArray {
					tags = append(tags, tag[1])
					esc := url.QueryEscape(tag[1])
					sb.WriteString(fmt.Sprintf("<a class=\"tag\" href=\"?tag=%s\">%s</a>", esc, tag[1]))
					sb.WriteString(" ")
				}
				sb.WriteString("</p>")
			}
			continue
		}
		//convert url to a tag
		ss := linkReg.ReplaceAllString(s, "<a href=\"$1\">$1</a>")
		sb.WriteString(fmt.Sprintf("<p>%s</p>", ss))
	}
	return strings.ReplaceAll(sb.String(), "\"", "\\\""), tags
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
	Attachments []Attachment
	ReleaseTime time.Time
}

type Attachment struct {
	FileType     string `json:"type"`
	FileMimeType string `json:"originType"`
	FileURL      string `json:"url"`
}
