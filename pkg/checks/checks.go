package checks

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"healthcheck/pkg/db"
	"healthcheck/pkg/models"
	"healthcheck/pkg/notifications"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func DoChecksWithInterval(logger *logrus.Entry, connDb *sql.DB, configurations *models.Configuration) error {
	// each 20 second we launch our healthcheck
	startTime, err := strconv.ParseInt(os.Getenv("StartTime"), 10, 64)
	if err != nil {
		logrus.Errorf("failed strconv.ParseInt: %s", err)
		return err
	}

	for range time.NewTicker(time.Second * time.Duration(startTime)).C {
		for _, value := range configurations.Configs {
			logger.Infof("checks url: %s", value.Url)
			results, err := CheckStatusCodeAndText(connDb, value.Url)
			if err != nil {
				logrus.Errorf("failed CheckStatusCodeAndText: %s", err)
				return err
			}
			if results.Checks != nil {
				logrus.Infof("%s: %s, %v\n", results.Url, results.Text, results.Checks)
			}
			logrus.Infof("%s: %s\n", results.Url, results.Text)

			logger.Infof("finished checking url: %s", value.Url)
		}
	}
	return nil
}

func CheckStatusCodeAndText(connDb *sql.DB, url string) (*models.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorf("failed http.Get: %s", err)
		return nil, err
	}
	response := &models.Response{
		Url:        url,
		StatusCode: resp.StatusCode,
		Text:       "ok",
	}
	ok := IsOk(resp)
	if !ok {
		response.Text = "fail"
		response.Checks = []string{"status_code", "text"}
	}

	savedResults, updatedResults, err := db.SendResultsOfChecksToDb(connDb, response)
	if err != nil {
		logrus.Errorf("failed SendResultsOfChecksToDb: %s", err)
		return nil, err
	}

	if savedResults != updatedResults {
		_, err = notifications.NotifyHttpBin(os.Getenv("HttpUrl"))
		if err != nil {
			logrus.Errorf("failed NotifyHttpBin(): %s", err)
			return nil, err
		}
	}

	return response, nil
}

func IsOk(resp *http.Response) bool {
	ok := strings.Contains(resp.Status, "OK")
	if !ok || resp.StatusCode != 200 {
		return false
	}
	return true
}
