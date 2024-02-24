package zeroaccount

type Metadata struct {
	userID           string
	profileID        string
	isWebhookRequest bool
}

func (m *Metadata) IsWebhookRequest() bool {
	return m.isWebhookRequest
}

func (m *Metadata) UserID() string {
	return m.userID
}

func (m *Metadata) ProfileID() string {
	return m.profileID
}
