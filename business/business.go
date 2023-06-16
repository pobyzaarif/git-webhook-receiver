package business

type InternalContext struct {
	TrackerID string
}

type WebhookProvider string

const (
	Github    WebhookProvider = "github"
	Gitlab    WebhookProvider = "gitlab"
	Bitbucket WebhookProvider = "bitbucket"
)

func NewInternalContext(trackerID string) InternalContext {
	return InternalContext{
		TrackerID: trackerID,
	}
}
