package utils

import (
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func AppendStrings(strings ...string) string {
	buf := make([]byte, 0)
	for _, str := range strings {
		buf = append(buf, str...)
	}
	return string(buf)
}

func IsValidId(id string) bool {
	r := regexp.MustCompile(`(?m)^[0-9a-zA-Z-]+$`)
	return r.MatchString(id)
}

func CreateUuid() string {
	uuid := uuid.NewV4().String()
	return uuid
}

func CreateApiKey() string {
	uuid := uuid.NewV4().String()
	return strings.Replace(uuid, "-", "", -1)
}
