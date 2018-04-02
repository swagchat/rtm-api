package metrics

import (
	"time"

	"github.com/swagchat/rtm-api/logging"
	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap/zapcore"
)

type stdoutProvider struct{}

func (sp *stdoutProvider) Run() {
	c := utils.Config()
	exec(func() {
		l, _ := time.LoadLocation("Etc/GMT")
		nowTime := time.Unix(time.Now().Unix(), 0).In(l)
		m := makeMetrics(nowTime)
		sb := utils.NewStringBuilder()
		mStr := sb.PrintStruct("config", m)
		logging.Log(zapcore.InfoLevel, &logging.AppLog{
			Kind:    "metrics",
			Message: mStr,
		})
	}, c.Metrics.Interval)
}
