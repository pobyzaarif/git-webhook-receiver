package controller

import (
	"net/http"
	"strings"

	v10 "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pobyzaarif/git-webhook-receiver/business"
	webhookNotif "github.com/pobyzaarif/git-webhook-receiver/business/webhookNotif"
	goLogger "github.com/pobyzaarif/go-logger/logger"
)

type Controller struct {
	webhookNotifService webhookNotif.Service
	validator           *v10.Validate
}

func NewController(webhookNotifService webhookNotif.Service) *Controller {
	return &Controller{
		webhookNotifService,
		v10.New(),
	}
}

var logger = goLogger.NewLog("CONTROLLER")

type (
	githubPayload struct {
		Ref        string `json:"ref"` // Branch Name
		Repository struct {
			Name string `json:"name"` // Repo Name
		} `json:"repository"`
		Pusher struct {
			Name string `json:"name"` // Author
		} `json:"pusher"`
	}
)

func (controller *Controller) GithubController(c echo.Context) error {
	trackerID, _ := c.Get("tracker_id").(string)
	logger.SetTrackerID(trackerID)
	ic := business.NewInternalContext(trackerID)

	request := new(githubPayload)
	if err := c.Bind(request); err != nil {
		logger.Error(err.Error(), err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": http.StatusText(http.StatusBadRequest)})
	}

	branchName := ""
	branchNameSplit := strings.Split(request.Ref, "/")
	if len(branchNameSplit) > 1 {
		branchName = branchNameSplit[len(branchNameSplit)-1] // take last word of path refs/heads/main >> main
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": http.StatusText(http.StatusInternalServerError)})
	}
	spec := webhookNotif.WebhookNotifSpec{
		WebhookProvider: business.Github,
		RepoName:        request.Repository.Name,
		BranchName:      branchName,
		Author:          request.Pusher.Name,
	}
	err := controller.webhookNotifService.Notif(ic, spec)
	if err != nil {
		logger.Error(err.Error(), err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
	}

	logger.Info("ok")
	return c.JSON(http.StatusOK, map[string]interface{}{"message": http.StatusText(http.StatusOK)})
}
