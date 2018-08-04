package metrics

import (
	"github.com/swagchat/rtm-api/config"
)

type stdoutProvider struct{}

func (sp *stdoutProvider) Run() {
	c := config.Config()
	exec(func() {
		// TODO
	}, c.Metrics.Interval)
}
