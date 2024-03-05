package zeroaccount

type Metadata struct {
	userID           string
	profileID        string
	isWebhookRequest bool
}

func (m *Metadata) IsWebhookRequest() bool {
	return m.isWebhookRequest
}

// UserID returns the user ID associated with the request.
// Deprecated: it is not actually deprecated. The method is for internal use only.
func (m *Metadata) UserID() string {
	return m.userID
}

func (m *Metadata) ProfileID() string {
	return m.profileID
}
