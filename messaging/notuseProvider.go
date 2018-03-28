package messaging

type NotuseProvider struct{}

func (provider *NotuseProvider) Subscribe() {
	// Do not process anything
}

func (provider *NotuseProvider) Unsubscribe() {
	// Do not process anything
}
