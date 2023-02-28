package notifications

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func NotifyHttpBin(url string) (string, error) {
	jsonData, err := json.Marshal("sent failed url")
	if err != nil {
		logrus.Errorf("failed json.Marshal: %s", err)
		return "", err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Errorf("failed http.Post: %s", err)
		return "", err
	}
	if resp.StatusCode != 200 {
		logrus.Error("resp.StatusCode not 200")
		return "", err
	}
	return "OK", err
}
