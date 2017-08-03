package utils

const (
	API_VERSION = "v0"
)

var (
	Realtime RealtimeSetting
	Que      QueSetting
)

type RealtimeSetting struct {
	Port string
}

type QueSetting struct {
	Host    string
	Port    string
	Topic   string
	Channel string
}
