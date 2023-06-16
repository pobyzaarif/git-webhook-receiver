package webhooknotif

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/pobyzaarif/git-webhook-receiver/business"
	"github.com/pobyzaarif/git-webhook-receiver/config"

	validator "github.com/go-playground/validator/v10"
	goLogger "github.com/pobyzaarif/go-logger/logger"
)

var logger = goLogger.NewLog("SERVICE")

type WebhookNotifSpec struct {
	WebhookProvider business.WebhookProvider `validate:"required"`
	RepoName        string                   `validate:"required"`
	BranchName      string                   `validate:"required"`
	Author          string                   `validate:"required"`
}

type (
	service struct {
		WebhookSetting *config.WebhookSetting
		validate       *validator.Validate
	}

	Service interface {
		Notif(ic business.InternalContext, spec WebhookNotifSpec) (err error)
	}
)

func NewService(WebhookSetting *config.WebhookSetting) Service {
	return &service{
		WebhookSetting,
		validator.New(),
	}
}

func (s *service) Notif(ic business.InternalContext, spec WebhookNotifSpec) (err error) {
	logger.SetTrackerID(ic.TrackerID)

	// Validate Spec
	if err = s.validate.Struct(spec); err != nil {
		logger.ErrorWithData("err invalid spec", map[string]interface{}{"spec": spec}, err)
		return
	}

	// Choose mapping
	var mappingUsed []config.Mapping
	switch spec.WebhookProvider {
	default:
		err = errors.New("not implemented yet")
		return
	case business.Github:
		mappingUsed = s.WebhookSetting.Github.Mapping
	case business.Gitlab:
		mappingUsed = s.WebhookSetting.Gitlab.Mapping
	case business.Bitbucket:
		mappingUsed = s.WebhookSetting.Bitbucket.Mapping
	}

	// Command to execute the shell command
	command := ""
	for _, mapping := range mappingUsed {
		if mapping.RepoName == spec.RepoName &&
			mapping.BranchName == spec.BranchName &&
			mapping.Command != "" {
			command = mapping.Command
			break
		}
	}

	if command == "" {
		logger.InfoWithData("ignored", map[string]interface{}{"spec": spec})
		return nil
	}

	// Execute the shell command
	output, err := exec.Command("sh", "-c", command).CombinedOutput()
	if err != nil {
		logger.ErrorWithData("err exec command", map[string]interface{}{"command": command}, err)
		return
	}

	res := fmt.Sprintf("Command output:\n%s", output)
	logger.InfoWithData("ok", map[string]interface{}{"output": res})

	return nil
}
