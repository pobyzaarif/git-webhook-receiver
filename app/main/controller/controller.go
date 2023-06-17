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

	gitlabPayload struct {
		Ref        string `json:"ref"` // Branch Name
		Repository struct {
			Name string `json:"name"` // Repo Name
		} `json:"repository"`
		UserUsername string `json:"user_username"` // Author
	}

	bitbucketPayload struct {
		Push struct {
			Changes []struct {
				New struct {
					Name string `json:"name"` // Branch Name
				} `json:"new"`
			} `json:"changes"`
		} `json:"push"`
		Repository struct {
			FullName string `json:"full_name"` // Repo Name
		} `json:"repository"`
		Actor struct {
			DisplayName string `json:"display_name"` // Author
		} `json:"actor"`
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
	return controller.controllerHelper(c, ic, spec)
}

func (controller *Controller) GitlabController(c echo.Context) error {
	trackerID, _ := c.Get("tracker_id").(string)
	logger.SetTrackerID(trackerID)
	ic := business.NewInternalContext(trackerID)

	request := new(gitlabPayload)
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
		WebhookProvider: business.Gitlab,
		RepoName:        request.Repository.Name,
		BranchName:      branchName,
		Author:          request.UserUsername,
	}
	return controller.controllerHelper(c, ic, spec)
}

func (controller *Controller) BitbucketController(c echo.Context) error {
	trackerID, _ := c.Get("tracker_id").(string)
	logger.SetTrackerID(trackerID)
	ic := business.NewInternalContext(trackerID)

	request := new(bitbucketPayload)
	if err := c.Bind(request); err != nil {
		logger.Error(err.Error(), err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": http.StatusText(http.StatusBadRequest)})
	}

	repoName := ""
	repoNameSplit := strings.Split(request.Repository.FullName, "/")
	if len(repoNameSplit) > 1 {
		repoName = repoNameSplit[len(repoNameSplit)-1] // take last word of path refs/heads/main >> main
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": http.StatusText(http.StatusInternalServerError)})
	}
	spec := webhookNotif.WebhookNotifSpec{
		WebhookProvider: business.Bitbucket,
		RepoName:        repoName,
		BranchName:      request.Push.Changes[0].New.Name,
		Author:          request.Actor.DisplayName,
	}
	return controller.controllerHelper(c, ic, spec)
}

func (controller *Controller) controllerHelper(c echo.Context, ic business.InternalContext, spec webhookNotif.WebhookNotifSpec) error {
	err := controller.webhookNotifService.Notif(ic, spec)
	if err != nil {
		logger.Error(err.Error(), err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
	}

	logger.Info("ok")
	return c.JSON(http.StatusOK, map[string]interface{}{"message": http.StatusText(http.StatusOK)})
}
