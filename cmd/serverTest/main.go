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
	err := configs.InitConfig("/../../pkg/configs/appConfigs.env")
	if err != nil {
		panic(err)
	}

	d, err := db.GetDB()
	if err != nil {
		logrus.Error("couldn't instantiate db")
	}
	defer d.Close()

	err = db.MigrateDb("../../migrations")
	if err != nil {
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
