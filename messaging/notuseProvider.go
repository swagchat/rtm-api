package messaging

type notuseProvider struct{}

func (provider *notuseProvider) Subscribe() {
	// Do not process anything
}

func (provider *notuseProvider) Unsubscribe() {
	// Do not process anything
}
