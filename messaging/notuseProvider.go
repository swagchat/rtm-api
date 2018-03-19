package messaging

type NotUseProvider struct{}

func (provider NotUseProvider) Init() error {
	return nil
}

func (provider NotUseProvider) Subscribe() {
}

func (provider NotUseProvider) Unsubscribe() {
}
