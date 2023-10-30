package halo

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

const templateDir = "/app/templates"

func NewMoment(m Moment) {
	tplPath := fmt.Sprintf("%s/halo-text.tpl", templateDir)
	if m.ImageURL != "" {
		tplPath = fmt.Sprintf("%s/halo-image.tpl", templateDir)
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
