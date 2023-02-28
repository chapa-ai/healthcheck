package controllers

import (
	"github.com/labstack/echo/v4"
	"healthcheck/pkg/models"
	"net/http"
	"strings"
)

func MakeRequest(c echo.Context) error {
	params := &models.Params{}
	err := c.Bind(params)
	if err != nil {
		return err
	}

	resp, err := http.Post(params.Url, "application/json", nil)
	if err != nil {
		return err
	}

	text := ""
	ok := strings.Contains(resp.Status, "OK")
	if !ok {
		text = "fail"
	}
	text = "ok"

	results := &models.Response{
		Url:        params.Url,
		StatusCode: resp.StatusCode,
		Text:       text,
	}

	return c.JSON(http.StatusOK, results)
}
