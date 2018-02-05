package utils

const (
	API_VERSION      = "v0"
	MAX_MESSAGE_SIZE = 8192
)

var (
	Realtime RealtimeSetting
	Que      QueSetting
)

type RealtimeSetting struct {
	Port                    string
	IsDisplayConnectionInfo bool
}

type QueSetting struct {
	Port           string
	NsqlookupdHost string
	NsqlookupdPort string
	NsqdHost       string
	NsqdPort       string
	Topic          string
	Channel        string
}
