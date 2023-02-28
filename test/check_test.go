package test

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"healthcheck/pkg/checks"
	"healthcheck/pkg/configs"
	"healthcheck/pkg/db"
	"healthcheck/pkg/models"
	"testing"
)

var (
	url = "http://127.0.0.1:9997"
)

func TestDoChecksWithInterval(t *testing.T) {

	err := configs.InitConfig("/../configs/appConfigs.env")
	if err != nil {
		t.Fatalf("failed get configs: %s", err)
	}

	params := &models.Params{
		Url: "https://ya.ru",
	}

	positiveResponse, err := RequestWithPositiveResponse(params)
	if err != nil {
		t.Fatalf("failed RequestWithPositiveResponse: %v", err)
	}

	badResponse, err := RequestWithBadResponse(params)
	if err != nil {
		t.Fatalf("failed RequestWithBadResponse: %v", err)
	}

	if badResponse == positiveResponse {
		t.Fatal("should not be equal")
	}
}

func TestCheckForBadRequests(t *testing.T) {
	err := configs.InitConfig("/../configs/appConfigs.env")
	if err != nil {
		t.Fatalf("failed get configs: %s", err)
	}
	connDb, err := db.GetDB()
	if err != nil {
		t.Fatalf("failed GetDB: %s", err)
	}
	maxAcceptableQuantityOfBadRequests := 2

	params := &models.Params{
		Url: "http://ya.ru",
	}
	counts := []int{}

	for i := 0; i <= 4; i++ {
		response, err := checks.CheckStatusCodeAndText(context.Background(), connDb, params.Url)
		if err != nil {
			t.Fatalf("failed checks.CheckStatusCodeAndText: %s", err)
		}
		counts = append(counts, response.StatusCode)
	}

	freq := make(map[int]int)
	for _, num := range counts {
		freq[num] = freq[num] + 1
	}

	for key, value := range freq {
		if key != 200 && value >= maxAcceptableQuantityOfBadRequests {
			t.Fatal("too many bad requests")
		}
	}
}

func TestSendResultsOfChecksToDb(t *testing.T) {
	err := configs.InitConfig("/../configs/appConfigs.env")
	if err != nil {
		t.Fatalf("failed get configs: %s", err)
	}
	connDb, err := db.GetDB()
	if err != nil {
		t.Fatalf("failed GetDB: %s", err)
	}
	resultsWithFail := &models.Response{
		Url:        "http://ya.ru",
		StatusCode: 500,
		Text:       "fail",
	}
	resultsWithOk := &models.Response{
		Url:        "http://ya.ru",
		StatusCode: 200,
		Text:       "ok",
	}

	_, updatedResultsWithFail, err := db.SendResultsOfChecksToDb(context.Background(), connDb, resultsWithFail)
	if err != nil {
		t.Fatalf("failed SendResultsOfChecksToDb: %s", err)
	}

	_, updatedResultsWithOk, err := db.SendResultsOfChecksToDb(context.Background(), connDb, resultsWithOk)
	if err != nil {
		t.Fatalf("failed SendResultsOfChecksToDb: %s", err)
	}

	if updatedResultsWithFail == updatedResultsWithOk {
		t.Fatal("should not be equal")
	}
}

func RequestWithPositiveResponse(params *models.Params) (*models.Response, error) {
	output := &models.Response{}

	url := fmt.Sprintf("%s/results", url)

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		SetResult(output).
		Post(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("status code wrong. status: %d. body: %s", resp.StatusCode(), resp.String())
	}

	return output, nil
}

func RequestWithBadResponse(params *models.Params) (*models.Response, error) {
	output := models.Response{}
	url := fmt.Sprintf("%s/results", url)

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(params.Url).
		SetResult(&output).
		Post(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 200 {
		return nil, fmt.Errorf("status code wrong. status: %d. body: %s", resp.StatusCode(), resp.String())
	}

	return &output, nil
}
