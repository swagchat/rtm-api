package metrics

import (
	"fmt"
	"time"

	logger "github.com/betchi/zapper"
	"github.com/kylelemons/godebug/pretty"
	"github.com/swagchat/rtm-api/config"
)

type stdoutProvider struct{}

func (sp *stdoutProvider) Run() {
	c := config.Config()
	exec(func() {
		l, _ := time.LoadLocation("Etc/GMT")
		nowTime := time.Unix(time.Now().Unix(), 0).In(l)
		m := makeMetrics(nowTime)

		compact := &pretty.Config{
			Compact: true,
		}
		logger.Info(fmt.Sprintf("Metrics: %s", compact.Sprint(m)))
	}, c.Metrics.Interval)
}
