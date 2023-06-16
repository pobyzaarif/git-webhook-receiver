package router

import (
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
	"github.com/pobyzaarif/git-webhook-receiver/app/main/controller"
	"github.com/pobyzaarif/git-webhook-receiver/config"
)

var apiVersion = "v1"

func RegisterPath(
	e *echo.Echo,
	conf *config.Config,
	controller *controller.Controller,
) {
	// Middleware for basic authentication
	basicAuthMiddleware := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Validate username and password here
		// You can replace the validation logic with your own implementation
		if username == conf.AppSetting.AppBasicAuthUsername && password == conf.AppSetting.AppBasicAuthPassword {
			return true, nil
		}
		return false, nil
	})

	e.GET("", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{"message": "hello world"})
	})

	notification := e.Group(apiVersion+"/webhook_notification", basicAuthMiddleware)
	notification.POST("/github", controller.GithubController)
}
