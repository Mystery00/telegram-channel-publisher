package halo

import (
	"bytes"
	"net/http"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

func NewMoment(m Moment) {
	tplPath := "templates/halo-text.tpl"
	if m.ImageURL != "" {
		tplPath = "templates/halo-image.tpl"
	}
	t := template.Must(template.ParseGlob(tplPath))
	var tmplBytes bytes.Buffer
	err := t.Execute(&tmplBytes, m)
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

type Moment struct {
	Content       string
	ImageMimeType string
	ImageURL      string
	ReleaseTime   time.Time
}
