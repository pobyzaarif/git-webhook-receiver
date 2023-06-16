package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pobyzaarif/git-webhook-receiver/app/main/controller"
	"github.com/pobyzaarif/git-webhook-receiver/app/main/router"
	webhookNotifBusiness "github.com/pobyzaarif/git-webhook-receiver/business/webhookNotif"
	"github.com/pobyzaarif/git-webhook-receiver/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	goLoggerAppName "github.com/pobyzaarif/go-logger/appname"
	goLogger "github.com/pobyzaarif/go-logger/logger"
	goLoggerEchoMiddlerware "github.com/pobyzaarif/go-logger/rest/framework/echo/v4/middleware"
)

var logger = goLogger.NewLog("MAIN")

func main() {
	conf := config.LoadConfig("./config.json")

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(goLoggerEchoMiddlerware.ServiceRequestTime)
	e.Use(goLoggerEchoMiddlerware.ServiceTrackerID)
	e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: goLoggerEchoMiddlerware.APILogHandler,
		Skipper: goLoggerEchoMiddlerware.DefaultSkipper,
	}))

	e.Use(goLoggerEchoMiddlerware.Recover())

	webhookNotifService := webhookNotifBusiness.NewService(&conf.WebhookSetting)
	ctrl := controller.NewController(webhookNotifService)

	router.RegisterPath(
		e,
		conf,
		ctrl,
	)

	address := "0.0.0.0:" + conf.AppSetting.AppPort
	go func() {
		if err := e.Start(address); err != http.ErrServerClosed {
			logger.Fatal("failed on http server " + conf.AppSetting.AppPort)
		}
	}()

	logger.SetTrackerID("main")
	logger.Info(goLoggerAppName.GetAPPName() + " service running in " + address)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// a timeout of 10 seconds to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("failed to shutting down echo server %v", err))
	} else {
		logger.Info("successfully shutting down echo server")
	}
}
