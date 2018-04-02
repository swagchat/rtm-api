package metrics

type notuseProvider struct{}

func (np *notuseProvider) Run() {
	// Do not process anything
}
