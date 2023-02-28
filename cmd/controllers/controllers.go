package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"healthcheck/pkg/db"
	"healthcheck/pkg/models"
	"net/http"
)

func GetFailedChecks(c echo.Context) error {
	res, err := db.GetFailedChecks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(fmt.Sprintf("failed GetFailedChecks: %v", err)))
	}

	return c.JSON(http.StatusOK, res)
}
