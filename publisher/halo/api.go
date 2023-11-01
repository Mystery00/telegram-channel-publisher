package halo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"telegram-channel-publisher/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	host  = viper.GetString(config.HaloHost)
	token = viper.GetString(config.HaloToken)
)

func getReqHost() string {
	if strings.HasSuffix(host, "/") {
		return host[:len(host)-1]
	} else {
		return host
	}
}

func request(uri, method string, body []byte) ([]byte, error) {
	return requestWithContentType(uri, method, "application/json", bytes.NewReader(body))
}

func requestWithContentType(uri, method, contentType string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", getReqHost(), uri)
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	logrus.Debugf("request halo: %s %s", method, url)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != http.StatusOK {
		logrus.Debugf("request halo error: %d, response: %s", resp.StatusCode, string(response))
		return response, fmt.Errorf("request halo error: %d", resp.StatusCode)
	}
	return response, nil
}
