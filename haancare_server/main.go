package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var err error

	InitLogger()

	err = InitConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to init config: %+v", err))
		os.Exit(1)
	}

	err = InitDB()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to init db: %+v", err))
		os.Exit(1)
	}
	defer db.Close()

	err = GenTables()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to gen db tables: %+v", err))
		os.Exit(1)
	}

	alpsToken, err = GetLatestAlpsToken()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get alps token: %+v", err))
		os.Exit(1)
	}

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// register cron jobs
	err = InitCronJobs()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to start cron jobs: %w", err))
		os.Exit(1)
	}

	InitRouter(e)
	e.Logger.Fatal(e.Start(":8080"))
}