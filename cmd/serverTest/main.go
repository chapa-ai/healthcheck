package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"healthcheck/cmd/serverTest/controllers"
	"healthcheck/pkg/configs"
	"healthcheck/pkg/db"
	"os"
)

func main() {
	err := configs.InitConfig("/../../configs/appConfigs.env")
	if err != nil {
		logrus.Error("couldn't get configs")
		panic(err)
	}

	d, err := db.GetDB()
	if err != nil {
		logrus.Error("couldn't instantiate db")
		return
	}
	defer func() {
		err = d.Close()
		if err != nil {
			logrus.Errorf("failed closing conn of db: %s", err)
			return
		}
	}()

	err = db.MigrateDb("../../migrations")
	if err != nil {
		logrus.Errorf("failed MigrateDb: %s", err)
		panic(err)
	}
	logrus.Info("migrations implemented")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/results", controllers.MakeRequest)

	err = e.Start(os.Getenv("ServerTestPort"))
	if err != nil {
		panic(err)
	}

}
