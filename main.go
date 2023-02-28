package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"healthcheck/cmd/controllers"
	"healthcheck/pkg/checks"
	"healthcheck/pkg/configs"
	"healthcheck/pkg/db"
	"os"
)

func main() {
	logger := logrus.WithFields(logrus.Fields{})
	logger.Info("program started")

	err := configs.InitConfig("/pkg/configs/appConfigs.env")
	if err != nil {
		panic(err)
	}

	d, err := db.GetDB()
	if err != nil {
		logrus.Errorf("couldn't instantiate db: %s", err)
		return
	}
	defer d.Close()

	err = db.MigrateDb("migrations")
	if err != nil {
		panic(err)
	}
	logger.Info("migrations implemented")

	configurations, err := configs.ReadUrlConfigs("/pkg/configs/configs.json")
	if err != nil {
		logrus.Errorf("failed read configs: %s", err)
		return
	}

	go func(logger *logrus.Entry) {
		err = checks.DoChecksWithInterval(logger, d, configurations)
		if err != nil {
			logrus.Errorf("failed checks.DoChecksWithInterval: %s", err)
		}
	}(logger)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/results", controllers.GetFailedChecks)

	err = e.Start(os.Getenv("AppPort"))
	if err != nil {
		panic(err)
	}

}
